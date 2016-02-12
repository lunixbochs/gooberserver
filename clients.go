package main

type Clients map[Client]bool

func (clients Clients) Contains(c Client) bool {
	_, ok := clients[c]
	return ok
}

func (clients Clients) Add(c Client) bool {
	if clients.Contains(c) {
		return false
	}
	clients[c] = true
	return true
}

func (clients Clients) Remove(c Client) bool {
	if !clients.Contains(c) {
		return false
	}
	delete(clients, c)
	return true
}
