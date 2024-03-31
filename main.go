package main

import (
	filereader "fast-rw/file_reader"
	"fmt"
	"time"
)

func main() {
	tReader := filereader.FileReader{}
	output := make(chan string)
	done := make(chan bool)

	go func() {
		tReader.ReadTokens("file_reader/reader.go", output, done)
	}()

	var tokens []string
	completed := false
	for {
		select {
		case token := <-output:
			tokens = append(tokens, token)
			fmt.Println(token)
			// sleep for 100ms to simulate processing
			time.Sleep(100 * time.Millisecond)
		case <-done:
			close(output)
			close(done)
			completed = true
		}
		if completed {
			break
		}
	}

	fmt.Println(len(tokens), " File Tokens Processed")
}
