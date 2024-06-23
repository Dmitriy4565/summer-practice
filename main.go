package main

import (
	"log"
	"net/http"

	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

const MaxConnections = 10

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type UserData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	defer db.Close()

}

var (
	db              *sql.DB
	userData        = make(map[string]string)
	mySigningKey    = []byte("secret")
	regexName       = regexp.MustCompile(`^[A-Za-z]+$`)
	regexPassword   = regexp.MustCompile(`^[A-Za-z0-9!@#$%&*]{8,32}$`)
	clients         = make(map[*websocket.Conn]bool)
	broadcast       = make(chan Message)
	LensConnections = 0
	upgrader        = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
)

func init() {
	var err error
	connStr := "user=postgres password=12345678 dbname=productdb sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func createToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Разрешен только метод POST", http.StatusMethodNotAllowed)
		return
	}

	var userData UserData
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := validate(userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (name, password) VALUES ($1, $2)", userData.Name, userData.Password)
	if err != nil {
		http.Error(w, "Ошибка при сохранении", http.StatusInternalServerError)
		return
	}

	tokenString, err := createToken(userData.Name)
	if err != nil {
		http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	fmt.Fprintln(w, "Пользователь зарегестрирован")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Разрешен только метод POST", http.StatusMethodNotAllowed)
		return
	}

	var userData UserData
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE name = $1", userData.Name).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка при запросе к БД", http.StatusInternalServerError)
		}
		return
	}

	if storedPassword != userData.Password {
		http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
		return
	}

	tokenString, err := createToken(userData.Name)
	if err != nil {
		http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	fmt.Fprintln(w, "Пользователь успешно вошел ")
}

func isStrongPassword(password string) bool {
	if len(password) < 8 || len(password) > 32 {
		return false
	}
	if !regexPassword.MatchString(password) {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%&*]`).MatchString(password)

	return hasUpper && hasLower && hasSpecial
}

func isNameValid(name string) bool {
	return regexName.MatchString(name)
}

func validate(userData UserData) error {
	if !isNameValid(userData.Name) {
		return fmt.Errorf("поле name должно содержать только буквы")
	}
	if !isStrongPassword(userData.Password) {
		return fmt.Errorf("пароль должен быть от 8 до 32 символов и содержать буквы, цифры и специальные символы")
	}
	return nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	if LensConnections >= MaxConnections {
		http.Error(w, "Превышено количество соединений", http.StatusTooManyRequests)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка при апгрейде WebSocket соединения: %v", err)
		return
	}

	LensConnections++
	clients[ws] = true
	defer func() {
		ws.Close()
		LensConnections--
		delete(clients, ws)
	}()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Ошибка при чтении WebSocket сообщения: %v", err)
			} else {
				log.Printf("WebSocket соединение закрыто нормально")
			}
			break
		}

		broadcast <- msg
	}
}
func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Ошибка при отправке WebSocket сообщения: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
