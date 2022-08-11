package main

import (
	"chat-server/src/message"
	"encoding/gob"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":6760")
	if err != nil {
		log.Println("connect error: ", err)
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
	for i := 0; i < 10; i++ {
		m.Option = message.Quit
		m.Payload = i
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

//go func(readChan chan bool, writeChan chan bool) {
//	for i := 0; i < 10; i++ {
//		log.Println("send", i)
//		writeChan <- true
//		<-readChan
//	}
//	close(writeChan)
//}(readChan, writeChan)
//
//go func(readChan chan bool, writeChan chan bool, exitChan chan bool) {
//	for {
//		v, ok := <-writeChan
//		if v && ok {
//			readChan <- true
//		} else {
//			break
//		}
//	}
//	exitChan <- true
//	close(readChan)
//	close(exitChan)
//}(readChan, writeChan, exitChan)
