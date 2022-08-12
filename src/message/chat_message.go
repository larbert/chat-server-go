package message

const (
	Error   byte = 0
	ACK     byte = 1
	Join    byte = 2
	Search  byte = 3
	Connect byte = 4
	Quit    byte = 5
)

type ChatMessage struct {
	Option  byte
	Payload interface{}
}
