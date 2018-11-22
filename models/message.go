package models

// should
type Message struct {
	// DateCreated string       `json:"dateCreated"`
	SessionID  string `json:"sessionID"`
	IssueTitle string `json:"issueTitle"`
	// Users [] connectedUsers
	StoryPoints int          `json:"storyPoints,omitempty"` // represents story point values for a giv
	Estimations []Estimation `json:"estimations"`
}
