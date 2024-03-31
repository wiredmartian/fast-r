package main

import (
	filereader "fast-r/file_reader"
	"fmt"
	"log"
	"time"
)

var (
	workers = 2
)

func main() {
	tokenReader := filereader.FileReader{}
	output := make(chan string)
	input := make(chan string)
	done := make(chan bool)
	var tokens []string

	// read all tokens from file
	tokens, err := tokenReader.ReadAllTokens("file.txt")
	if err != nil {
		log.Fatal(err)
	}

	// create workers to process the tokens
	for i := 0; i < workers; i++ {
		go tokenReader.Worker(input, output, done)
	}

	// send tokens to workers
	go func() {
		for i := 0; i < len(tokens); i++ {
			input <- tokens[i]
		}
		close(done)
	}()

	// print the tokens
	for {
		finished := false
		select {
		case t := <-output:
			// simulate processing
			time.Sleep(100 * time.Millisecond)
			fmt.Println("Line : ", t)
		case <-done:
			finished = true
		}
		if finished {
			break
		}
	}
	fmt.Println(len(tokens), " File Tokens Processed")
}
