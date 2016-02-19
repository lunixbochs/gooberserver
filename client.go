package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Client interface {
	Acl() AclFlags
	User() *User
	Nick() string
	Send(msg string)
}

type TCPClient struct {
	conn   net.Conn
	user   User
	acl    AclFlags
	server *Server
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
	c.conn.Write([]byte(msg + "\n"))
}

func (c *TCPClient) Close() {
	c.conn.Close()
}

func (c *TCPClient) Pump(remClient chan *TCPClient) {
	protocol := &Protocol{c, c.server}
	defer func() {
		protocol.removeClient()
		c.Close()
		remClient <- c
	}()
	protocol.addClient()
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		protocol.Handle(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error listening to client:", err)
	}
}
