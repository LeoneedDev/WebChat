package main

import (
	"log"
	"net/http"
	"github.com/LeoneedDev/WebChat/handler"
	
)

func main() {
	s := http.Server{
		Addr: ":8080",
	}
	
	http.Handle("/", handler.)

	}())
	if err := s.ListenAndServe(); err != nil {
		log.Fatal("On start server not work ",err)
	}

}