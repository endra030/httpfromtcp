package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	//fmt.Println("I hope I get the job!")
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}
	c := getLinesChannel(file)
	for line := range c {
		fmt.Printf("read: %s\n", line)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	retChan := make(chan string)
	go func() {
		defer f.Close()
		currentLine := ""
		for {
			data := make([]byte, 8)
			_, err := f.Read(data)
			if err == io.EOF {
				if currentLine != "" {
					retChan <- currentLine
					//fmt.Printf("read: %s\n", currentLine)

				}
				close(retChan)
				break

			}
			currentLine = currentLine + string(data)
			parts := strings.Split(currentLine, "\n")
			if len(parts) == 2 {
				//fmt.Printf("read: %s\n", parts[0])
				retChan <- parts[0]
				currentLine = "" + parts[1]
			}

		}
	}()
	return retChan

}
