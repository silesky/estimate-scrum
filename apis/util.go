package apis

import (
	"bytes"
	"encoding/json"
	"estimate/daos"
	"io"
	"log"
	"net/http"
)

// setup CORS
func SetupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Println(err)
	}
}

type Query struct {
	SessionID string
	AdminID   string
}

func GetQuery(r *http.Request) Query {
	sessionID := r.URL.Query().Get("id")
	adminID := r.URL.Query().Get("adminID")
	return Query{
		SessionID: sessionID,
		AdminID:   adminID,
	}
}

func IsAdmin(r *http.Request) bool {
	q := GetQuery(r)
	session, _ := daos.GetSession(q.SessionID)
	return q.AdminID == session.AdminID
}
