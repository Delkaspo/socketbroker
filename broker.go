package socketbroker

import (
	"log"
)

type Client interface {
	Send(interface{})
}

type SocketBroker struct {
	Name string

	bcast       chan interface{}
	subscribe   chan Client
	unsubscribe chan Client
	clients     map[Client]bool
}

func (sb *SocketBroker) Subscribe(c Client) {
	sb.subscribe <- c
}

func (sb *SocketBroker) Unsubscribe(c Client) {
	sb.unsubscribe <- c
}

func (sb *SocketBroker) Broadcast(d interface{}) {
	sb.bcast <- d
}

func (sb *SocketBroker) Run() {
	for {
		select {
		case c := <-sb.subscribe:
			sb.clients[c] = true
			log.Printf("New subscriber ! (live %d)\n", len(sb.clients))
		case c := <-sb.bcast:
			for v, _ := range sb.clients {
				v.Send(c)
			}
		case c := <-sb.unsubscribe:
			delete(sb.clients, c)
			log.Printf("A client left :( (live %d)\n", len(sb.clients))
		}
	}
}

func New(name string) *SocketBroker {
	return &SocketBroker{
		Name:        name,
		bcast:       make(chan interface{}),
		subscribe:   make(chan Client),
		unsubscribe: make(chan Client),
		clients:     make(map[Client]bool),
	}
}

