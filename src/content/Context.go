package content

import (
	"chat-server/src/message"
	"net"
)

type Context struct {
	Conn net.TCPConn
	Req  message.ChatMessage
	Res  message.ChatMessage
}
