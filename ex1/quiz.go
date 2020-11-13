package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type question struct {
	question string
	answer   string
}

type quiz struct {
	answered  int
	score     int
	questions []question
}

func loadQuiz(filePath string) *quiz {
	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Error openings CSV file", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	var questions []question
	for {
		quiz, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if quiz[1] == "" {
			continue
		}

		q := question{quiz[0], quiz[1]}
		questions = append(questions, q)
	}

	return &quiz{
		answered:  0,
		score:     0,
		questions: questions,
	}
}

func (q *quiz) run(limit int) {
	scanner := bufio.NewScanner(os.Stdin)

	timer := time.NewTimer(time.Duration(limit) * time.Second)
	for _, question := range q.questions {
		ansCh := make(chan string)
		fmt.Printf("%s = ?: ", question.question)
		go func(scanner *bufio.Scanner) {
			scanner.Scan()
			ans := scanner.Text()
			ansCh <- ans
		}(scanner)

		select {
		case <-timer.C:
			return
		case ans := <-ansCh:
			if ans == question.answer {
				q.score++
			}
			q.answered++
		}
	}
}

func (q quiz) report() {
	fmt.Printf("Score: %d out of %d", q.score, len(q.questions))
}

func main() {
	filePath := flag.String("file", "problems.csv", "CSV File that contains quiz.")
	limit := flag.Int("limit", 30, "Time limit to answer quiz.")

	flag.Parse()

	quiz := loadQuiz(*filePath)
	quiz.run(*limit)
	quiz.report()
}
