package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		stdOut, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		_, err = udpConn.Write([]byte(stdOut))
		if err != nil {
			log.Fatal(err)
		}

	}

}
