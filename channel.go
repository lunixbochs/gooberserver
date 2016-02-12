package main

import (
	"fmt"
)

type Channel struct {
	Name     string
	Owner    UserId
	Clients  Clients
	Admins   []UserId
	Ban      []UserId
	Allow    []UserId
	MuteList []UserId
	Topic    string
	Key      string

	Autokick  bool
	ChanServ  bool
	Antispam  bool
	Censor    bool
	Antishock bool
	Log       bool

	eventSend chan<- interface{}
	eventRecv <-chan interface{}
	OnClose   func(c *Channel)
}

func NewChannel(name string) *Channel {
	c := &Channel{Name: name}
	events := make(chan interface{})
	c.eventSend = events
	c.eventRecv = events
	c.Clients = make(Clients)
	go c.Pump()
	return c
}

// Channel events
type ChanJoin struct {
	Client Client
}

type ChanLeave struct {
	Client Client
}

type ChanMsg struct {
	Client Client
	Msg    string
	Ex     bool
}

// main event loop for Channel
func (c *Channel) Pump() {
	defer close(c.eventSend)
	for {
		event := <-c.eventRecv
		switch e := event.(type) {
		case ChanJoin:
			if c.Clients.Add(e.Client) {
				c.Broadcast("JOINED %s %s", c.Name, e.Client.Nick())
			}
		case ChanLeave:
			if c.Clients.Remove(e.Client) {
				c.Broadcast("LEFT %s %s", c.Name, e.Client.Nick())
			}
		case ChanMsg:
			if !c.Clients.Contains(e.Client) {
				continue
			}
			say := "SAY"
			if e.Ex {
				say = "SAYEX"
			}
			c.Broadcast("%s %s %s %s", say, c.Name, e.Client.Nick(), e.Msg)
		}
	}
}

// send a raw server message to all users in Channel
func (c *Channel) Broadcast(msg string, args ...string) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args)
	}
	for client, _ := range c.Clients {
		client.Send(msg)
	}
}
