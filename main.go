package main

import (
	"encoding/json"
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

// grabs the message from the client and sends it to the channel
func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println(r)

	sessionID := getQuery(r).sessionID

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

func parseBodyToUserMessageEstimationDTO(r *http.Request) (models.UserMessageEstimation, error) {
	var estimation models.UserMessageEstimation
	err := json.NewDecoder(r.Body).Decode(&estimation)
	return estimation, err
}

// handles POSTS to /api/estimations
func handleUpdateOrAddEstimation(w http.ResponseWriter, req *http.Request) {
	estimationDto, err := parseBodyToUserMessageEstimationDTO(req)
	if err != nil {
		log.Printf("Body parsing error for estimation:", "error: %v", err)
	}

	dbSaveError := daos.UpdateUserEstimation(estimationDto)
	if dbSaveError != nil {
		http.Error(w, "Could not save estimation to DB.", http.StatusNotFound) // 4
	}

	// encode struct as json and write to stream
	w.Header().Set("Content-Type", "application/json")
	// order is important
	w.WriteHeader(http.StatusCreated)
	// serialize the struct as JSON (again)
	fmt.Println("Estimation!", estimationDto)
	json.NewEncoder(w).Encode(estimationDto)

}

// create a new session with /new -- should probably be a POST
func handleCreateNewSession(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	session, err := daos.CreateNewSession()
	fmt.Println(session.ID)
	if err != nil {
		log.Printf("%v", err)
		// https://golang.org/pkg/net/http/#pkg-constants
		http.Error(w, "Unable to retrieve session.", http.StatusInternalServerError) // 500
		return
	}
	// encode struct as json and write to stream
	w.Header().Set("Content-Type", "application/json")
	// order is important
	w.WriteHeader(http.StatusCreated)
	// serialize the struct as JSON (again)
	fmt.Println("SESSION!", session)
	json.NewEncoder(w).Encode(session)
}

type Query struct {
	sessionID string
	adminID   string
}

func getQuery(req *http.Request) Query {
	sessionID := req.URL.Query().Get("id")
	adminID := req.URL.Query().Get("adminID")
	return Query{
		sessionID: sessionID,
		adminID:   adminID,
	}
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
func isAdmin(req *http.Request) bool {
	q := getQuery(req)
	session, _ := daos.GetSession(q.sessionID)
	return q.adminID == session.AdminID
}

// for requesting a specific session
func handleRequestSession(w http.ResponseWriter, req *http.Request) {
	q := getQuery(req)
	session, err := daos.GetSession(q.sessionID)
	if err != nil {
		log.Printf(q.sessionID)
		http.Error(w, "No session with that Session ID found.", http.StatusNotFound) // 404
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := session.GetSessionResponse(q.adminID)
	fmt.Printf("%+v\n", response)
	json.NewEncoder(w).Encode(response)
}

func parseBodyToSession(r *http.Request) (models.Session, error) {
	var session models.Session
	err := json.NewDecoder(r.Body).Decode(&session)
	return session, err
}

func handleUpdateSession(w http.ResponseWriter, req *http.Request) {
	q := getQuery(req)
	session, err := parseBodyToSession(req)
	if err != nil {
		http.Error(w, "Cannot parse req body.", http.StatusInternalServerError) // 404
		return
	}
	if !isAdmin(req) {
		http.Error(w, "Unauthorized user.", http.StatusUnauthorized)
		return
	}
	daos.UpdateSession(q.sessionID, session)
	// encode struct as json and write to stream
	w.Header().Set("Content-Type", "application/json")
	// order is important
	w.WriteHeader(http.StatusOK)
	// serialize the struct as JSON (again)
	json.NewEncoder(w).Encode("updated successfully.")

}

// setup CORS
func setupResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func handleSession(w http.ResponseWriter, req *http.Request) {
	setupResponse(&w)
	switch method := req.Method; method {
	case "GET":
		handleRequestSession(w, req)
	case "POST":
		// this should just return the new sessionID and adminKey, and client can use that information to navigate to the new url.
		handleCreateNewSession(w, req)
	case "PUT":
		handleUpdateSession(w, req)
	}
}

func handleUserEstimation(w http.ResponseWriter, req *http.Request) {
	setupResponse(&w)
	switch method := req.Method; method {
	case "POST":
		// this should just return the new sessionID and adminKey, and client can use that information to navigate to the new url.
		handleUpdateOrAddEstimation(w, req)
	}
}

func main() {

	db.Init()
	clientSessions = db.WsStore.Users
	handler := http.FileServer(http.Dir("./client"))
	// https://rickyanto.com/understanding-go-standard-http-libraries-servemux-handler-handle-and-handlefunc/
	http.Handle("/", handler)
	http.HandleFunc("/api/session", handleSession)
	http.HandleFunc("/api/estimation", handleUserEstimation)
	http.HandleFunc("/ws", handleConnections)
	// go handleIncomingMessages()
	go deliverMessages()

	port := ":3333"
	fmt.Println("Listenining on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
