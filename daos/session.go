package daos

import (
	"encoding/json"
	"errors"
	"estimate/db"
	"estimate/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/imdario/mergo"
)

// gets a session from the database
func GetSession(sessionID string) (models.Session, error) {
	// create pointer to a session
	sessionPtr := &models.Session{}

	// get the json data from db (will come back as byte[])
	sessionDbRes, err := db.Get(sessionID)

	// deserialize byte[] to a pointer to a Session.
	json.Unmarshal(sessionDbRes, sessionPtr)

	// dereference pointer and return struct (otherwise json response fields will come back empty.
	return *sessionPtr, err
}

// creates empty session
func GetDefaultSession() models.Session {
	return models.Session{
		DateCreated: time.Now().UTC().String(),
		ID:          uuid.New().String(), // public ID will allow others to connect to this session. Will be used as the redis key.
		AdminID:     uuid.New().String(),
		StoryPoints: []int{1, 2, 3},
		Issues:      []models.Issue{},
	}
}

func CreateNewSession() (models.Session, error) {
	defaultSession := GetDefaultSession()
	err := SaveSession(defaultSession.ID, defaultSession)
	return defaultSession, err
}

// if adminId is correct, update session
func UpdateSession(sessionID string, newData models.Session) error {
	session, err := GetSession(sessionID)
	if err != nil {
		panic(err)
	}
	log.Println(newData.AdminID, session.AdminID)
	if newData.AdminID != session.AdminID {
		return errors.New("Auth error")
	}
	mergo.Merge(&newData, session)
	SaveSession(session.ID, newData)
	return nil
}

func SaveSession(sessionID string, session models.Session) error {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		fmt.Printf("%v", "Could not convert to json.")
	}
	db.Set(sessionID, []byte(sessionJSON))
	return err
}

func GetDefaultIssue() models.Issue {
	return models.Issue{
		IssueTitle:  "",
		IssueID:     uuid.New().String(),
		Estimations: make(map[string]int),
	}
}

func CreateNewIssue(sessionID string) error {
	session, err := GetSession(sessionID)
	if err != nil {
		return errors.New("cannot find session.")
	}
	session.Issues = append(session.Issues, GetDefaultIssue())
	return SaveSession(sessionID, session)
}

// for updating the estimates when there's a new message
func getIssueByID(sessionID string, issueID string) (models.Issue, error) {
	session, err := GetSession(sessionID)
	issue := models.Issue{}
	for i := range session.Issues {
		eachIssue := &session.Issues[i]
		if eachIssue.IssueID == issueID {
			issue = *eachIssue
		}
	}
	return issue, err
}

// for updating the estimates when there's a new message
func updateIssue(sessionID string, newIssue models.Issue) error {
	session, err := GetSession(sessionID)
	for i := range session.Issues {
		eachIssue := &session.Issues[i]
		if eachIssue.IssueID == newIssue.IssueID {
			eachIssue = &newIssue
		}
	}
	SaveSession(sessionID, session)
	return err
}

// should take username, estimate (number), and issueID
func UpdateEstimations(sessionID string, issueID string, username string, estimate int) error {
	issue, err := getIssueByID(sessionID, issueID)
	if issue.IssueID == "" {
		return errors.New("Unable to create estimation. No issue with issue ID found")
	}
	if err != nil {
		panic(err)
	}
	issue.Estimations[username] = estimate
	return updateIssue(sessionID, issue)
}
