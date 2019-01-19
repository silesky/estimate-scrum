package handlers

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

// helpers

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// grabs the message from the client and sends it to the channel
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	apis.SetupCORS(&w)
	sessionID := apis.GetQuery(r).SessionID

	fmt.Println(sessionID)
	// upgrade initial GET to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	if db.ClientSessions[sessionID] == nil {
		db.ClientSessions[sessionID] = map[*websocket.Conn]bool{}
	}
	db.ClientSessions[sessionID][ws] = true
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
			delete(db.ClientSessions[sessionID], ws)
			break
		}
		db.WsStore.Broadcast <- msg
	}

	defer ws.Close()
}

// handles POSTS to /api/estimations
func HandleUpdateOrAddEstimation(w http.ResponseWriter, r *http.Request) {
	var e models.UserMessageEstimation
	err := apis.Decode(r, &e)
	if err != nil {
		apis.Respond(w, r, http.StatusInternalServerError, "Unable to parse estimation.")
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
func HandleCreateNewSession(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	session, err := daos.CreateNewSession()
	fmt.Println(session.ID)
	if err != nil {
		apis.Respond(w, r, http.StatusInternalServerError, "Unable to retrieve session")
		return
	}
	apis.Respond(w, r, http.StatusCreated, session)
}

func sendDataToClient(sessionID string, data models.SessionResponse) {
	for client := range db.ClientSessions[sessionID] {
		err := client.WriteJSON(data)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(db.ClientSessions[sessionID], client)
		}
	}
}

func DeliverMessages() {
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
			log.Println()
			panic("error pub/sub, delivery has stopped")
		default:
			log.Println("DEFAULT CASE")
			log.Println(db.PSC.Receive())
		}
	}
}

// for ruesting a specific session
func HandleRequestSession(w http.ResponseWriter, r *http.Request) {
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

func HandleUpdateSession(w http.ResponseWriter, r *http.Request) {
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

func HandleSession(w http.ResponseWriter, r *http.Request) {
	apis.SetupCORS(&w)
	switch method := r.Method; method {
	case "GET":
		HandleRequestSession(w, r)
	case "POST":
		// this should just return the new sessionID and adminKey, and client can use that information to navigate to the new url.
		HandleCreateNewSession(w, r)
	case "PUT":
		HandleUpdateSession(w, r)
	}
}

func HandleUserEstimation(w http.ResponseWriter, r *http.Request) {
	apis.SetupCORS(&w)
	switch method := r.Method; method {
	case "POST":
		// this should just return the new sessionID and adminKey, and client can use that information to navigate to the new url.
		HandleUpdateOrAddEstimation(w, r)
	}
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./client"))
}
