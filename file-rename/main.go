package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	entries, err := os.ReadDir("./test_files")

	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fmt.Println("Name", e.Name())
		fmt.Println("IsDir", e.IsDir())
	}
}
