package filereader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	SMALL_FILE  int64 = (1024 * 1024) * 10  // 10MB
	MEDIUM_FILE int64 = (1024 * 1024) * 100 // 100MB
	LARGE_FILE  int64 = (1024 * 1024) * 500 // 500MB
)

type Reader interface {
	// ReadAllTokens reads all tokens from a file and returns them as a slice of strings
	ReadAllTokens(path string) ([]string, error)
	// ReadTokens reads tokens from a file and sends them to a channel
	ReadTokens(path string, output chan string, done chan bool) error
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
		f.tokens = append(f.tokens, line)
	}
	return f.tokens, nil
}

func (f *FileReader) ReadTokens(path string, output chan string, done chan bool) error {
	var err error

	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()

	if err != nil {
		return err
	}

	workers := 1
	if info.Size() > SMALL_FILE && info.Size() < MEDIUM_FILE {
		workers = 2
	} else if info.Size() > MEDIUM_FILE && info.Size() < LARGE_FILE {
		workers = 3
	} else if info.Size() > LARGE_FILE {
		workers = 6
	}

	fmt.Println("Number of workers: ", workers)

	reader := bufio.NewReader(file)
	// read file line by line
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// done reading file
				done <- true
				break
			}
			break
		}
		output <- line
	}
	return nil
}

func (f *FileReader) ProcessFileTokens(workers int, output chan string, done chan bool) error {
	for i := 0; i < workers; i++ {
		go f.worker(output, output, done)
	}
	return nil
}

func (f *FileReader) worker(input chan string, output chan string, done chan bool) {
	for {
		select {
		case data := <-input:
			output <- data
		case <-done:
			return
		}
	}
}
