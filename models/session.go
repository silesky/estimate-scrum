package models

import validation "github.com/go-ozzo/ozzo-validation"

type Session struct {
	AdminID       string  `json:"adminID,omitempty"`
	DateCreated   string  `json:"dateCreated"`
	ID            string  `json:"ID"`
	StoryPoints   []int   `json:"storyPoints"`
	Issues        []Issue `json:"issues"`
	SelectedIssue string  `json:"selectedIssue"`
}

func (s Session) OK() error {
	return validation.ValidateStruct(
		validation.Field(s.ID, validation.Required),
	)
}

type SessionResponse struct {
	Session Session `json:"session"`
	IsAdmin bool    `json:"isAdmin"`
}

// get sanitized session reponse (this request will get sent to every websocket client.)
func (s Session) GetSessionResponse(adminID string) SessionResponse {
	s.AdminID = "" // sanitize session
	return SessionResponse{
		Session: s,
		IsAdmin: adminID == s.AdminID,
	}
}
