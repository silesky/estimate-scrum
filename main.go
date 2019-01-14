package main

import (
	"encoding/json"
	"estimate/apis"
	"estimate/daos"
	"estimate/db"
	"estimate/models"
	"fmt"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
)

// globals
var (
	clientSessions db.WSUserMap
)

var broadcast = make(chan models.UserMessageEstimation)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// helpers

// grabs the message from the client and sends it to the channel
func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	apis.SetupCORS(&w)
	sessionID := apis.GetQuery(r).SessionID

	fmt.Println(sessionID)
	// upgrade initial GET to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	if clientSessions[sessionID] == nil {
		clientSessions[sessionID] = map[*websocket.Conn]bool{}
	}
	clientSessions[sessionID][ws] = true
	for {
		// add type switches so I can handle both AdminMessages and UserMessages
		var msg models.UserMessageEstimation
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			// TODO: this is wrong
			delete(clientSessions[sessionID], ws)
			break
		}
		broadcast <- msg
	}

	defer ws.Close()
}

func sendDataToClient(sessionID string, data models.SessionResponse) {
	for client := range clientSessions[sessionID] {
		err := client.WriteJSON(data)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(clientSessions[sessionID], client)
		}
	}
}

// handles POSTS to /api/estimations
func handleUpdateOrAddEstimation(w http.ResponseWriter, r *http.Request) {
	var e models.UserMessageEstimation
	err := apis.Decode(r, &e)
	if err != nil {
		apis.Respond(w, r, http.StatusInternalServerError, err)
		return
	}

	dbSaveError := daos.UpdateUserEstimation(e)
	if dbSaveError != nil {
		apis.Respond(w, r, http.StatusNotFound, "Could not save estimation to DB.")
		return
	}
	apis.Respond(w, r, http.StatusCreated, e)
}

// create a new session with /new -- should probably be a POST
func handleCreateNewSession(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	session, err := daos.CreateNewSession()
	fmt.Println(session.ID)
	if err != nil {
		apis.Respond(w, r, http.StatusInternalServerError, "Unable to retrieve session")
		return
	}
	apis.Respond(w, r, http.StatusCreated, session)
}

func deliverMessages() {
	log.Println("delivering messages.")
	for {

		switch v := db.PSC.Receive().(type) {
		case redis.PMessage:
			log.Printf("pmessage: %s: %s", v.Channel, v.Data)
			sessionID := string(v.Data)
			s, _ := daos.GetSession(sessionID)
			sendDataToClient(sessionID, s.GetSessionResponse(""))
			fmt.Printf("%+v\n", s)

		case redis.Message:
			log.Printf("message: %s: %s", v.Channel, v.Data)

		case redis.Subscription:

			log.Printf("subscription: %s: %s %d\n", v.Channel, v.Kind, v.Count)

		case error:
			log.Println("error pub/sub, delivery has stopped")

		default:
			log.Println("DEFAULT CASE")
			log.Println(db.PSC.Receive())
		}
	}
}

// for ruesting a specific session
func handleRequestSession(w http.ResponseWriter, r *http.Request) {
	q := apis.GetQuery(r)
	session, err := daos.GetSession(q.SessionID)
	if err != nil {
		apis.Respond(w, r, http.StatusNotFound, err)
		log.Printf(q.SessionID)
		return
	}
	data := session.GetSessionResponse(q.AdminID)
	apis.Respond(w, r, http.StatusOK, data)
}

func handleUpdateSession(w http.ResponseWriter, r *http.Request) {
	q := apis.GetQuery(r)
	var session models.Session
	err := apis.Decode(r, &session)
	if err != nil {
		apis.Respond(w, r, http.StatusInternalServerError, err)
		fmt.Println("cannot parse client JSON body.")
		return
	}
	if !apis.IsAdmin(r) {
		apis.Respond(w, r, http.StatusUnauthorized, "Unauthorized user.")
		return
	}
	daos.UpdateSession(q.SessionID, session)
	apis.Respond(w, r, http.StatusOK, "updated successfully.")

}

func handleSession(w http.ResponseWriter, r *http.Request) {
	apis.SetupCORS(&w)
	switch method := r.Method; method {
	case "GET":
		handleRequestSession(w, r)
	case "POST":
		// this should just return the new sessionID and adminKey, and client can use that information to navigate to the new url.
		handleCreateNewSession(w, r)
	case "PUT":
		handleUpdateSession(w, r)
	}
}

func handleUserEstimation(w http.ResponseWriter, r *http.Request) {
	apis.SetupCORS(&w)
	switch method := r.Method; method {
	case "POST":
		// this should just return the new sessionID and adminKey, and client can use that information to navigate to the new url.
		handleUpdateOrAddEstimation(w, r)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./client"))
}

func main() {

	db.Init()
	clientSessions = db.WsStore.Users

	// https://rickyanto.com/understanding-go-standard-http-libraries-servemux-handler-handle-and-handlefunc/
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/session", handleSession)
	http.HandleFunc("/api/estimation", handleUserEstimation)
	http.HandleFunc("/ws", handleConnections)
	// go handleIncomingMessages()
	go deliverMessages()

	port := ":3333"
	fmt.Println("Listenining on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
