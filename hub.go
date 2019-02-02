package main

import (
	"log"
	"strconv"
)

// BroadCast is room config
type BroadCast struct {
	to  string
	msg []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	// clients map[*Client]bool
	clients map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *BroadCast

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *BroadCast),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if _, ok := h.clients[client.room]; !ok {
				h.clients[client.room] = map[*Client]bool{}
			}
			h.clients[client.room][client] = true
			for cls := range h.clients[client.room] {
				select {
				case cls.send <- []byte(`{"total":` + strconv.Itoa(len(h.clients[client.room])) + `, "header":"open", "from":"` + client.id + `"}`):
				default:
					close(cls.send)
					delete(h.clients[cls.room], cls)
				}
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client.room][client]; ok {
				delete(h.clients[client.room], client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients[message.to] {
				log.Println(client)
				select {
				case client.send <- message.msg:
				default:
					close(client.send)
					delete(h.clients[client.room], client)
				}
			}
		}
	}
}
