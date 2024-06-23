package main

import (
	"sync"
	"ThirdTask/generator"
	"ThirdTask/handlers"
	"ThirdTask/utils"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	ch := make(chan int)

	go generator.GenerateNumbers(ch, &wg)
	go handlers.DivisibleByTwo(ch, &wg)
	go handlers.DivisibleByThree(ch, &wg)
	go handlers.DivisibleByFour(ch, &wg)

	wg.Wait()

	utils.WriteToFile("divisibleByTwo", handlers.DivisibleByTwoResults)
	utils.WriteToFile("divisibleByThree", handlers.DivisibleByThreeResults)
	utils.WriteToFile("divisibleByFour", handlers.DivisibleByFourResults)
}
