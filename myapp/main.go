package main

import (
	"myapp/generator"
	"myapp/handlers"
	"myapp/utils"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	ch := make(chan int)

	// Генерация чисел
	go generator.GenerateNumbers(ch, &wg)
	// Обработка кратных 2
	go handlers.DivisibleByTwo(ch, &wg)
	// Обработка кратных 3
	go handlers.DivisibleByThree(ch, &wg)
	// Обработка кратных 4
	go handlers.DivisibleByFour(ch, &wg)

	wg.Wait()

	// Запись результатов в файл
	utils.WriteToFile("divisibleByTwo", handlers.DivisibleByTwoResults)
	utils.WriteToFile("divisibleByThree", handlers.DivisibleByThreeResults)
	utils.WriteToFile("divisibleByFour", handlers.DivisibleByFourResults)
}
