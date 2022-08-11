package server

import (
	"chat-server/src/message"
	"encoding/gob"
	"log"
	"net"
)

type HandlerFunc func(interface{})

type Server struct {
}

func (s *Server) Start() {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":6760")
	tcp, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}
	defer func(tcp *net.TCPListener) {
		err := tcp.Close()
		if err != nil {
			return
		}
	}(tcp)
	for {
		acceptTCP, err := tcp.AcceptTCP()
		if err != nil {
			return
		}
		go handleTcp(acceptTCP)
	}
}

func handleTcp(conn *net.TCPConn) {
	defer func(conn *net.TCPConn) {
		err := conn.Close()
		if err != nil {
			log.Println("close error: ", err)
		}
	}(conn)
	for {
		// 创建解码器
		dec := gob.NewDecoder(conn)
		enc := gob.NewEncoder(conn)
		// 创建消息结构体
		m := &message.ChatMessage{}
		// 解码
		err := dec.Decode(m)
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				log.Println("connect useless")
				break
			} else {
				log.Println("decode error: ", err)
			}
		} else {
			// 处理收到的消息
			log.Println(m)
		}
		// 返回响应消息
		err = enc.Encode(&message.ChatMessage{
			Option:  message.ACK,
			Payload: nil,
		})
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				log.Println("connect useless")
				break
			} else {
				log.Println("encode error: ", err)
			}
		}
	}
}
