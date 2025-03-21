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
		
		for {
			

			}
		}
	}()
	return retChan

}
