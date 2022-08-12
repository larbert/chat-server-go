package main

import (
	"bufio"
	"chat-server/src/message"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	input := bufio.NewReader(os.Stdin)
	fmt.Printf("输入TCP目的地: ")
	tcpAddr, _ := input.ReadString('\n')
	tcpAddr = strings.TrimSpace(tcpAddr)
	conn, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		log.Println("connect error: ", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("close error: ", err)
		}
	}(conn)

	exitChan := make(chan bool)
	readChan := make(chan bool)
	writeChan := make(chan bool)
	go handleWrite(conn, readChan, writeChan)
	go handleRead(conn, readChan, writeChan, exitChan)

	for ok := range exitChan {
		log.Println(ok)
	}
	log.Println("主线程结束！")
}

func handleWrite(conn net.Conn, readChan chan bool, writeChan chan bool) {
	m := &message.ChatMessage{}
	input := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		messageString, _ := input.ReadString('\n')
		messageString = strings.TrimSpace(messageString)
		if messageString == "bye" {
			break
		}
		m.Option = message.Quit
		m.Payload = messageString
		enc := gob.NewEncoder(conn)
		err := enc.Encode(m)
		if err != nil {
			log.Println("Error to send message because of ", err)
		}
		log.Println("send ", m)
		writeChan <- true
		<-readChan
	}
	close(writeChan)
}

func handleRead(conn net.Conn, readChan chan bool, writeChan chan bool, exitChan chan bool) {
	m := &message.ChatMessage{}
	for {
		flag, ok := <-writeChan
		if !flag || !ok {
			break
		}
		dec := gob.NewDecoder(conn)
		err := dec.Decode(m)
		if err != nil {
			log.Println(err)
		}
		log.Println(m)
		readChan <- true
	}
	exitChan <- true
	close(readChan)
	close(exitChan)
}
