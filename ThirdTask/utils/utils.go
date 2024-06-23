package utils

import (
	"fmt"
	"os"
)

func WriteToFile(filename string, results []string) {
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
