package main

import (
	"fmt"
	// "strings"
)

type Protocol struct {
	Client Client
	Server *Server
}

func (p *Protocol) addClient() {
	p.Client.Send("TASServer 0.0 0.0 8201 1")
}

func (p *Protocol) removeClient() {
}

func (p *Protocol) Handle(msg string) {
	fmt.Println("protocol handler:", msg)
}
