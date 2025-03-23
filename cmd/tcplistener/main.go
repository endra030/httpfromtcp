package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	//fmt.Println("I hope I get the job!")
	l, err := net.Listen("tcp4", "127.0.0.1:42069")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		con, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		fmt.Println("connection has been accepted..")
		c := getLinesChannel(con)
		for line := range c {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("connection has been closed..")

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
