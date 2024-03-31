package filereader

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Reader interface {
	// ReadAllTokens reads all tokens from a file and returns them as a slice of strings
	ReadTokens(path string) ([]string, error)
}

type FileReader struct {
	tokens []string
}

func (f *FileReader) ReadAllTokens(path string) ([]string, error) {
	var err error

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// read file line by line
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// done reading file
				break
			}
			break
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		f.tokens = append(f.tokens, line)
	}
	return f.tokens, nil
}

func (f *FileReader) Worker(input <-chan string, output chan<- string, done <-chan bool) {
	for {
		select {
		case token := <-input:
			output <- "**" + token
		case <-done:
			return
		}
	}
}
