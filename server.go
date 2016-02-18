package main

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	Clients  map[Client]bool
	Channels map[string]*Channel
	Battles  map[int]*Battle
}

func NewServer() *Server {
	return &Server{
		Clients:  make(map[Client]bool),
		Channels: make(map[string]*Channel),
		Battles:  make(map[int]*Battle),
	}
}

func (s *Server) ListenAndServe(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	addClient := make(chan *TCPClient)
	remClient := make(chan *TCPClient)
	go func() {
		for {
			conn, err := ln.Accept()
			// TODO: check for max open files
			if err != nil {
				fmt.Println("error accepting connection:", err)
				os.Exit(1)
			}
			addClient <- &TCPClient{
				conn: conn, server: s,
				acl:     AclFlags{AclEveryone, AclUnauthed},
				onClose: remClient,
			}
		}
	}()
	for {
		select {
		case client := <-addClient:
			s.Clients[client] = true
		case client := <-remClient:
			delete(s.Clients, client)
		}
	}
}

func ListenAndServe(addr string) error {
	server := NewServer()
	return server.ListenAndServe(addr)
}

func main() {
	err := ListenAndServe("localhost:8200")
	if err != nil {
		fmt.Printf("Error serving: %s", err)
		os.Exit(1)
	}
}
