package main

import (
	"encoding/json"
	"estimate/models"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// globals
var clients = make(map[*websocket.Conn]bool) // create in-memory global map to keep track of client data

var broadcast = make(chan models.Estimation)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println(r)

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
		log.Printf("msg: %v", msg)
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

func handleCreateNewSession(w http.ResponseWriter, r *http.Request) {
	// create new uuid in redis, if successful, send response
	res := models.Session{
		DateCreated: time.Now().String(),
		ID:          uuid.New().String(), // public ID will allow others to connect to this session. Will be used as the redis key.
		AdminID:     uuid.New().String(),
		IssueTitle:  "",
		StoryPoints: []int{}, // initialize empty slice
	}
	json.NewEncoder(w).Encode(res)
}

func main() {
	handler := http.FileServer(http.Dir("./client"))
	// https://rickyanto.com/understanding-go-standard-http-libraries-servemux-handler-handle-and-handlefunc/
	http.Handle("/", handler)                 // handler is an instance of a ServeMux struct, not a fn.
	http.HandleFunc("/ws", handleConnections) // this is our normal fn

	// this should return a session id
	http.HandleFunc("/new", handleCreateNewSession)

	go handleEstimations()
	fmt.Println("Listing on 3333.")
	log.Fatal(http.ListenAndServe(":3333", nil))
}
