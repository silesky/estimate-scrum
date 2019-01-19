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

// Session that is sanitized of private data
func (s Session) Public() Session {
	s.AdminID = "" // sanitize session
	return s
}

// Session API response for admins and non-admins alike
func (s Session) Response(adminID string) SessionResponse {
	return SessionResponse{
		Session: s.Public(),
		IsAdmin: adminID == s.AdminID,
	}
}
