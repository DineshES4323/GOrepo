package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func ProcessLogs(inputFiles []string, outputFile string) error {

	errorChan := make(chan string)

	var a sync.WaitGroup

	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)

		file, err := os.Create(outputFile)
		if err != nil {
			log.Printf("Failed to create output file: %v\n", err)
			return
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		defer writer.Flush()

		for line := range errorChan {
			_, err := writer.WriteString(line + "\n")
			if err != nil {
				log.Printf("Failed to write to file: %v\n", err)
			}
		}
	}()

	for _, filename := range inputFiles {
		a.Add(1)
		go func(filePath string) {
			defer a.Done()

			file, err := os.Open(filePath)
			if err != nil {
				log.Printf("Failed to open file %s: %v\n", filePath, err)
				return
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "ERROR") {
					errorChan <- line
				}
			}
			if err := scanner.Err(); err != nil {
				log.Printf("Error reading file %s: %v\n", filePath, err)
			}
		}(filename)
	}

	go func() {
		a.Wait()
		close(errorChan)
	}()

	<-writerDone
	return nil
}

func main() {
	inputFiles := []string{"server1.log", "server2.log", "server3.log"}
	err := ProcessLogs(inputFiles, "errors.log")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Error log extraction completed.")
}
