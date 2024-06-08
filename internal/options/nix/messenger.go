package nix

import "fmt"

type Messenger struct{}

func NewMessenger() Messenger {
	return Messenger{}
}

func (Messenger) Send(msg string) {
	fmt.Println(msg)
}
