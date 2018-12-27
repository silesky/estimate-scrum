package models

// should
type UserMessageEstimation struct {
	Username        string `json:"username"`
	SessionID       string `json:"sessionID"`
	IssueID         string `json:"issueID"`
	EstimationValue int    `json:"estimationValue"`
}

/*
{
 username: "seth",
 sessionID: 456,
 issueID: 123,
 estimationValue: 123,
}
*/
type AdminMessage struct {
	// DateCreated string       `json:"dateCreated"`
	SessionID  string `json:"sessionID"`
	IssueTitle string `json:"issueTitle"`
	// Users [] connectedUsers
	StoryPoints int          `json:"storyPoints,omitempty"` // represents story point values for a giv
	Estimations []Estimation `json:"estimations"`

	Username string `json:"username"`
	Estimate int    `json:"estimate"`
	IssueID  string `json:"issueID"`
}
