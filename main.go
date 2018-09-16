package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) //
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi! there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Listing on 8080.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
