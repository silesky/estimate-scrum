package main

import (
	"encoding/json"
	"estimate/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// globals
var clients = make(map[*websocket.Conn]bool) // create in-memory global map to keep track of client data

var broadcast = make(chan models.Estimation)
var upgrader = websocket.Upgrader{}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// upgrade initial GET to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	clients[ws] = true

	for {
		var msg models.Estimation
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}

	defer ws.Close()
}

// goroutine that runs in a new thread
func handleEstimations() {
	for {
		// grab the next msg from the broadcast chanell
		msg := <-broadcast
		log.Printf("msg", msg)
		// iterate over each client
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	handler := http.FileServer(http.Dir("./client"))
	// https://rickyanto.com/understanding-go-standard-http-libraries-servemux-handler-handle-and-handlefunc/
	http.Handle("/", handler)                 // handler is an instance of a ServeMux struct, not a fn.
	http.HandleFunc("/ws", handleConnections) // this is our normal fn
	go handleEstimations()
	fmt.Println("Listing on 3333.")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
