package models

import validation "github.com/go-ozzo/ozzo-validation"

// should
type UserMessageEstimation struct {
	Username        string `json:"username"`
	SessionID       string `json:"sessionID"`
	IssueID         string `json:"issueID"`
	EstimationValue int    `json:"estimationValue"`
}

func (u UserMessageEstimation) OK() error {
	return validation.ValidateStruct(
		validation.Field(u.IssueID, validation.Required),
	)
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
	StoryPoints   int          `json:"storyPoints,omitempty"` // represents story point values for a giv
	Estimations   []Estimation `json:"estimations"`
	SelectedIssue string       `json:"SelectedIssue"`
	FinalPoints   int          `json:"FinalPoints"`
	IssueID       string       `json:"issueID"`
}
