package commands

import "errors"

var (
	UnknownCommand = errors.New("Unknown Command")
)

type SendCommand struct {
	Message string
}

type NameCommand struct {
	Name string
}

type MessageCommand struct {
	Name    string
	Message string
}
