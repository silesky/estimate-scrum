package main

import (
	"estimate/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// globals
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)
var upgrader = websocket.Upgrader{}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi! there, I love %s!", r.URL.Path[1:])
}

func main() {
	handler := http.FileServer(http.Dir("./client"))
	http.Handle("/", handler)
	fmt.Println("Listing on 3333.")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
