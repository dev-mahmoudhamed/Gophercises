package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	filePath := "problems.csv"
	totalQuestions := 0
	rightAnswers := 0
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

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

	totalQuestions = len(records)

	for i, record := range records {
		question := record[0]
		answer := record[1]

		fmt.Printf("Problem #%d: %s = ", i, question)

		answerCh := make(chan string)
		go readUserAnswer(answerCh)

		timeout := time.After(time.Second * time.Duration(*limit))

		select {
		case answerInput := <-answerCh:
			if answerInput == answer {
				rightAnswers++
			}
		case <-timeout:
			fmt.Printf("\nYou scored %d out of %d.\n", rightAnswers, totalQuestions)
			return
		}

	}
	fmt.Printf("\nYou scored %d out of %d.\n", rightAnswers, totalQuestions)
}

func readUserAnswer(ch chan<- string) {
	var answerInput string
	fmt.Scan(&answerInput)
	ch <- answerInput
}
