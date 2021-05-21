package client

import (
	"io"
	"log"
	"net"

	"github.com/stnnnghm/go-chat-server/commands"
)

type TcpChatClient struct {
	conn      net.Conn
	cmdReader *commands.CommandReader
	cmdWriter *commands.CommandWriter
	name      string
	incoming  chan commands.MessageCommand
}

func NewClient() *TcpChatClient {
	return &TcpChatClient{
		incoming: make(chan commands.MessageCommand),
	}
}

func (c *TcpChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err == nil {
		c.conn = conn
	}

	c.cmdReader = commands.NewCommandReader(conn)
	c.cmdWriter = commands.NewCommandWriter(conn)

	return err
}

func (c *TcpChatClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Read error: %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case commands.MessageCommand:
				c.incoming <- v

			default:
				log.Printf("Unknown Command: %v", v)
			}
		}
	}
}

func (c *TcpChatClient) Close() {
	c.conn.Close()
}

func (c *TcpChatClient) Incoming() chan commands.MessageCommand {
	return c.incoming
}

func (c *TcpChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

func (c *TcpChatClient) SetName(name string) error {
	return c.Send(commands.NameCommand{Name: name})
}

func (c *TcpChatClient) SendMessage(message string) error {
	return c.Send(commands.SendCommand{
		Message: message,
	})
}
