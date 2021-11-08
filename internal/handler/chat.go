package handler

import (
	"net/http"
	"github.com/LeoneedDev/WebChat/internal/chat"
)

func ChatHandler (w http.ResponseWriter, r *http.Request) {
	//Reg new chat
	hub := chat.New()

	//Run new chat in gorutine
	go hub.Run()

	//Serve our chat
	chat.ServeWs(hub,w,r)
}