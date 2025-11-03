package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	filePath := "problems.csv"
	totalQuestions := 0
	rightAnswers := 0

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// totalQuestions = len(records)

	for i, record := range records {
		question := record[0]
		answer := record[1]
		totalQuestions++

		fmt.Printf("Problem #%d: %s = ", i, question)

		var answerInput string
		fmt.Scan(&answerInput)

		if answerInput == answer {
			rightAnswers++
		}
	}
	fmt.Printf("You scored %d out of %d.\n", rightAnswers, totalQuestions)
}
