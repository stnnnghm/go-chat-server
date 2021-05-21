package main

import "github.com/stnnnghm/go-chat-server/server"

func main() {
	s := server.NewServer()
	s.Listen(":3333")

	s.Start()
}
