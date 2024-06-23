package handlers

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	DivisibleByTwoResults   []string
	DivisibleByThreeResults []string
	DivisibleByFourResults  []string
	mutex                   sync.Mutex
)

// DivisibleByTwo обрабатывает числа, кратные 2
func DivisibleByTwo(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		if num%2 == 0 {
			processResult("divisibleByTwo", num)
		}
		time.Sleep(time.Duration(rand.Intn(21)+10) * time.Millisecond)
	}
}

// DivisibleByThree обрабатывает числа, кратные 3
func DivisibleByThree(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		if num%3 == 0 {
			processResult("divisibleByThree", num)
		}
		time.Sleep(time.Duration(rand.Intn(71)+30) * time.Millisecond)
	}
}

// DivisibleByFour обрабатывает числа, кратные 4
func DivisibleByFour(ch <-chan int, wg *sync.WaitGroup) {
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
		DivisibleByTwoResults = append(DivisibleByTwoResults, fmt.Sprintf("%d %s\n", num, handler))
	case "divisibleByThree":
		DivisibleByThreeResults = append(DivisibleByThreeResults, fmt.Sprintf("%d %s\n", num, handler))
	case "divisibleByFour":
		DivisibleByFourResults = append(DivisibleByFourResults, fmt.Sprintf("%d %s\n", num, handler))
	}
}
