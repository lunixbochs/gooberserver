package main

import (
	"fmt"
	"net"
)

type Client interface {
	Acl() AclFlags
	User() *User
	Nick() string
	Send(msg string)
}

type TCPClient struct {
	conn    net.Conn
	user    User
	acl     AclFlags
	server  *Server
	onClose chan *TCPClient
}

func (c *TCPClient) Acl() AclFlags {
	// TODO: where do channel/battle get mixed in?
	return c.acl
}

func (c *TCPClient) User() *User {
	return &c.user
}

func (c *TCPClient) Nick() string {
	return c.user.Nick
}

func (c *TCPClient) Send(msg string) {
	fmt.Printf("sending %s to %s\n", c.Nick, msg)
}
