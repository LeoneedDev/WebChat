package main

import (
	"log"
	"net/http"
	"github.com/LeoneedDev/WebChat/internal/chat"
	"github.com/LeoneedDev/WebChat/internal/handler"
)

func main() {
	//Reg new chat
	hub := chat.New()

	//Run new chat in gorutine
	go hub.Run()

	
	s := http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/", handler.HomeHandler)
	http.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request)  {
		chat.ServeWs(hub,w,r)
	})
	log.Println("server start at ", s.Addr)

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
		return
	}

}
