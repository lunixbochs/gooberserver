package main

import (
	"fmt"
	"reflect"
	"strings"
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

func (p *Protocol) Respond(msg string) {
	// TODO: add msg id responder
	p.Client.Send(msg)
}

func (p *Protocol) Handle(msg string) {
	split := strings.Split(msg, " ")
	// TODO: parse out msg ids
	cmd := strings.ToUpper(split[0])
	args := split[1:]

	method := reflect.ValueOf(p).MethodByName(cmd)
	if !method.IsValid() {
		// TODO: check if DENIED is the right response
		p.Respond(fmt.Sprintf("DENIED Command not found: %s", cmd))
		return
	}
	numIn := method.Type().NumIn()
	if numIn > len(args) {
		p.Respond(fmt.Sprintf("DENIED Not enough arguments to %s", cmd))
		return
	} else if numIn < len(args) {
		extra := args[numIn-1:]
		args = args[:numIn-1]
		args = append(args, strings.Join(extra, " "))
	}
	argVals := make([]reflect.Value, len(args))
	for i, v := range args {
		argVals[i] = reflect.ValueOf(v)
	}
	// TODO: errors?
	// TODO: per-battle/channel acls
	acl := p.Client.Acl()
	if err := acl.Check(cmd, 0, 0); err != nil {
		// TODO: better message
		p.Respond(fmt.Sprintf("DENIED %s", err))
		return
	}
	method.Call(argVals)
}

func (p *Protocol) LOGIN(user, pass string) {
	fmt.Println("LOGIN", user, pass)
}
