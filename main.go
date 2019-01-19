package main

import (
	"estimate/db"
	"estimate/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {

	db.Init()

	// https://rickyanto.com/understanding-go-standard-http-libraries-servemux-handler-handle-and-Handlefunc/
	http.HandleFunc("/", handlers.HandleIndex)
	http.HandleFunc("/api/session", handlers.HandleSession)
	http.HandleFunc("/api/estimation", handlers.HandleUserEstimation)
	http.HandleFunc("/ws", handlers.HandleConnections)
	// go handleIncomingMessages()
	go handlers.DeliverMessages()

	port := ":3333"
	fmt.Println("Listenining on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
