package server

import (
	"chat-server/src/content"
	"chat-server/src/message"
	"chat-server/src/user"
	"encoding/gob"
	"log"
	"net"
	"reflect"
)

type HandlerFunc func(*content.Context)

var handleFuncMap = make(map[byte]HandlerFunc, 100)

type Server struct {
}

func (s *Server) Start() {
	// 注册要传输的gob对象
	Register()
	// 注册每个消息类型对应的方法
	registerFuncMap()

	// 启动服务
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":6760")
	tcp, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return
	}
	// 结束前关闭监听
	defer func(tcp *net.TCPListener) {
		err := tcp.Close()
		if err != nil {
			return
		}
	}(tcp)
	// 开始监听
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
		m := message.ChatMessage{}
		// 解码
		err := dec.Decode(&m)
		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				log.Println("connect useless")
				break
			} else if err.Error() == "EOF" {
				log.Println("connect finish")
				break
			} else {
				log.Println("decode error: ", err)
			}
		} else {
			// 处理收到的消息
			log.Println(m)
			context := &content.Context{
				Conn: *conn,
				Req:  m,
				Res:  message.ChatMessage{},
			}
			handleFuncMap[context.Req.Option](context)
			// 返回响应消息
			err = enc.Encode(context.Res)
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
}

func connect(ctx *content.Context) {
	ctx.Res.Option = message.Error
	ctx.Res.Payload = nil
}

func join(ctx *content.Context) {
	rv := reflect.ValueOf(ctx.Req.Payload)
	u := rv.Interface().(user.User)
	user.Join(u)
	ctx.Res.Option = message.ACK
	ctx.Res.Payload = user.UserList
}

func defaultHandleFunc(ctx *content.Context) {
	ctx.Res.Option = message.Error
	ctx.Res.Payload = nil
}

func registerFuncMap() {
	for i := 0; i < 10; i++ {
		handleFuncMap[byte(i)] = defaultHandleFunc
	}
	handleFuncMap[message.Join] = join
	handleFuncMap[message.Connect] = connect
}
