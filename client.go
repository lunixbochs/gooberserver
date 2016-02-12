package main

import (
	"fmt"
)

type Client interface {
	Acl() AclFlags
	User() *User
	Nick() string
	Send(msg string)
}

type RealClient struct {
	user User
	acl  AclFlags
}

func (c *RealClient) Acl() AclFlags {
	// TODO: where do channel/battle get mixed in?
	return c.acl
}

func (c *RealClient) User() *User {
	return &c.user
}

func (c *RealClient) Nick() string {
	return c.user.Nick
}

func (c *RealClient) Send(msg string) {
	fmt.Printf("sending %s to %s\n", c.Nick, msg)
}
