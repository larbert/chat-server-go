package main

import "chat-server/src/server"

func main() {
	s := &server.Server{}
	s.Start()
}
