package chat

import "log"

type Hub struct {
	//Registered clients
	clients map[*Client]bool

	//Messages from client
	broadcast chan []byte

	//Register req from clients
	register chan *Client

	//Unreg req from clients
	unregister chan *Client
}

func New() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}


//Run our hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			} else {
				log.Println("Try delete non exist cleint: ",client)
			}
			
		case messages := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- messages:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}