package main

import (
	"encoding/json"
	"estimate/daos"
	"estimate/db"
	"estimate/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// globals
var clients = make(map[*websocket.Conn]bool) // create in-memory global map to keep track of client data

// make a map of clients connected to
var clientSessions = make(map[string]map[*websocket.Conn]bool)

var broadcast = make(chan models.Estimation)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// grabs the message from the client and sends it to the channel
func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println(r)

	sessionID := r.URL.Query().Get("id")
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
		var msg models.Estimation
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

// broadcasts a message to all the users
func handleEstimations() {
	for {
		// grab the next msg from the broadcast chanell
		estimation := <-broadcast
		sessionID := estimation.SessionID
		log.Printf("msg: %v", estimation)

		dbSaveError := daos.UpdateEstimations(sessionID, estimation.IssueID, estimation.Username, estimation.Estimate)
		if dbSaveError != nil {
			log.Printf("error: %v", dbSaveError)
		}
		// iterate over each client
		for client := range clientSessions[sessionID] {
			err := client.WriteJSON(estimation)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clientSessions[sessionID], client)
			}
		}
	}
}
func createIssue(issueTitle string, issueID string, estimations models.EstimationsMap) models.Issue {
	return models.Issue{
		IssueID:     issueID,
		IssueTitle:  issueTitle,
		Estimations: estimations,
	}
}

// create a new session with /new -- should probably be a POST
func handleCreateNewSession(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	session, err := daos.CreateNewSession()
	daos.CreateNewIssue(session.ID)
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
	json.NewEncoder(w).Encode(session)
}

func isAdmin(adminID string, session models.Session) bool {
	return adminID == session.AdminID
}

// for requesting a specific session
func handleRequestSession(w http.ResponseWriter, req *http.Request) {
	sessionID := req.URL.Query().Get("id")
	adminID := req.URL.Query().Get("adminId")
	session, err := daos.GetSession(sessionID)
	if err != nil {
		log.Printf(sessionID)
		http.Error(w, "No session with that Session ID found.", http.StatusNotFound) // 404
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := session.GetSessionResponse(adminID)
	json.NewEncoder(w).Encode(response)
}

func parseBodyToSession(r *http.Request) (models.Session, error) {
	var session models.Session
	err := json.NewDecoder(r.Body).Decode(&session)
	return session, err
}

func handleUpdateSession(w http.ResponseWriter, req *http.Request) {
	sessionID := req.URL.Query().Get("id")
	session, err := parseBodyToSession(req)
	if err != nil {
		http.Error(w, "Cannot parse req body.", http.StatusInternalServerError) // 404
		return
	}
	updateError := daos.UpdateSession(sessionID, session)
	if updateError != nil {
		http.Error(w, "Unauthorized user.", http.StatusUnauthorized)
		return
	}
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

func main() {
	db.Init()
	handler := http.FileServer(http.Dir("./client"))
	// https://rickyanto.com/understanding-go-standard-http-libraries-servemux-handler-handle-and-handlefunc/
	http.Handle("/", handler)                 // handler is an instance of a ServeMux struct, not a fn.
	http.HandleFunc("/ws", handleConnections) // this is our normal fn

	// this should return a session if available, or an error message
	http.HandleFunc("/api/session", handleSession)
	go handleEstimations()
	port := ":3333"
	fmt.Println("Listenining on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
