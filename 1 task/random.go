package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	divisibleByTwoResults   []string
	divisibleByThreeResults []string
	divisibleByFourResults  []string
	mutex                   sync.Mutex
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	wg.Add(3)

	ch := make(chan int)

	go generateNumbers(ch, &wg)
	go divisibleByTwo(ch, &wg)
	go divisibleByThree(ch, &wg)
	go divisibleByFour(ch, &wg)

	wg.Wait()

	writeToFile("divisibleFour", divisibleByFourResults)
	writeToFile("divisibleByThree", divisibleByThreeResults)
	writeToFile("divisibleByTwo", divisibleByTwoResults)
}

func generateNumbers(ch chan<- int, wg *sync.WaitGroup) {
	defer close(ch)
	defer wg.Done()
	for i := 0; i < 50; i++ {
		num := rand.Intn(100) + 1
		ch <- num
	}
}

func divisibleByTwo(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		if num%2 == 0 {
			processResult("divisibleByTwo", num)
		}
		time.Sleep(time.Duration(rand.Intn(21)+10) * time.Millisecond)
	}
}

func divisibleByThree(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		if num%3 == 0 {
			processResult("divisibleByThree", num)
		}
		time.Sleep(time.Duration(rand.Intn(71)+30) * time.Millisecond)
	}
}

func divisibleByFour(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		if num%4 == 0 {
			processResult("divisibleByFour", num)
		}
		time.Sleep(time.Duration(rand.Intn(101)+100) * time.Millisecond)
	}
}

func processResult(handler string, num int) {
	mutex.Lock()
	defer mutex.Unlock()
	switch handler {
	case "divisibleByTwo":
		divisibleByTwoResults = append(divisibleByTwoResults, fmt.Sprintf("%d %s\n", num, handler))
	case "divisibleByThree":
		divisibleByThreeResults = append(divisibleByThreeResults, fmt.Sprintf("%d %s\n", num, handler))
	case "divisibleByFour":
		divisibleByFourResults = append(divisibleByFourResults, fmt.Sprintf("%d %s\n", num, handler))
	}
}

func writeToFile(filename string, results []string) {
	file, err := os.OpenFile(filename+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for _, result := range results {
		if _, err := file.WriteString(result); err != nil {
			fmt.Println(err)
			return
		}
	}
}
