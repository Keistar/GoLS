package main

import (
	"fmt"
	"os"
)

func main() {
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	fmt.Println("--- file list ---")
	for _, file := range files {
		prefix := "[FILE]"
		if file.IsDir() {
			prefix = "[DIR] "
		}
		fmt.Printf("%s %s\n", prefix, file.Name())
	}
}
