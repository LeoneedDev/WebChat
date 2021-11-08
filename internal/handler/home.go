package handler

import (
	"html/template"
	"log"
	"net/http"
)

//Home page handler
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	file, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Println("home.html not opened ",err)
		http.Error(w,"Page not uploded",http.StatusNotImplemented)
		return
	}
	err = file.Execute(w, nil)
	if err != nil {
		log.Println("home.html not execute", err)
		http.Error(w,"Page not uploded",http.StatusNotImplemented)
		return
	}
}