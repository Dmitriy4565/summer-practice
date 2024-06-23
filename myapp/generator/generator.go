package generator

import (
	"math/rand"
	"sync"
	"time"
)

// GenerateNumbers генерирует случайные числа и отправляет их в канал
func GenerateNumbers(ch chan<- int, wg *sync.WaitGroup) {
	defer close(ch)
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {
		num := rand.Intn(100) + 1
		ch <- num
	}
}
