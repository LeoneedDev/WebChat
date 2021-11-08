package chat

import (
	// "bytes"
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)


var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

//Client is in middle between websck conn and the hub
type Client struct {
	Hub *Hub
	
	//Websock conn
	Conn *websocket.Conn

	//Buffered channel of outbound messages
	Send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func (c *Client)readMessage() {
	defer func(){
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(
		func(string) error {
			 c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil 
			},
		)
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		c.Hub.broadcast <- message
	}

}


func (c *Client)writeMessage() {
	tiker := time.NewTicker(pongWait)
	defer func(){
		tiker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <- c.Send:
			err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println("DeadLine on write ", err)
				return
			}
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte("Close chan"))
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)

			if err != nil {
				log.Println("Next mess dont send ",err)
				return
			}

			w.Write(message)

			n := len(c.Send)

			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				log.Println("On close writer", err)
				return
			}
		case <-tiker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error on serveWs ",err)
		return
	}

	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.register <- client

	go client.readMessage()
	go client.writeMessage()
}


