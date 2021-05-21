package server

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"

	"github.com/stnnnghm/go-chat-server/commands"
)

type client struct {
	conn   net.Conn
	name   string
	writer *commands.CommandWriter
}

type TcpChatServer struct {
	listener net.Listener
	clients  []*client
	mutex    *sync.Mutex
}

var (
	UnknownClient = errors.New("Unknown Client")
)

func NewServer() *TcpChatServer {
	return &TcpChatServer{
		mutex: &sync.Mutex{},
	}
}

func (s *TcpChatServer) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err == nil {
		s.listener = l
	}

	log.Printf("Listening on %v", addr)

	return err
}

func (s *TcpChatServer) Close() {
	s.listener.Close()
}

func (s *TcpChatServer) Start() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Print(err)
		}

		client := s.accept(conn)
		go s.serve(client)
	}
}

func (s *TcpChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		client.writer.Write(command)
	}

	return nil
}

func (s *TcpChatServer) Send(name string, command interface{}) error {
	for _, client := range s.clients {
		if client.name == name {
			return client.writer.Write(command)
		}
	}

	return UnknownClient
}

func (s *TcpChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &client{
		conn:   conn,
		writer: commands.NewCommandWriter(conn),
	}

	s.clients = append(s.clients, client)

	return client
}

func (s *TcpChatServer) remove(client *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// remove connections from clients slice
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}

	log.Printf("Closing connection from %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}

func (s *TcpChatServer) serve(client *client) {
	cmdReader := commands.NewCommandReader(client.conn)

	defer s.remove(client)

	for {
		cmd, err := cmdReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read error: %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case commands.SendCommand:
				go s.Broadcast(commands.MessageCommand{
					Message: v.Message,
					Name:    client.name,
				})

			case commands.NameCommand:
				client.name = v.Name
			}
		}
	}
}
