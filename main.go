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
var broadcast = make(chan models.Estimation)
var upgrader = websocket.Upgrader{}

func main() {
	handler := http.FileServer(http.Dir("./client"))
	http.Handle("/", handler)
	fmt.Println("Listing on 3333.")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
