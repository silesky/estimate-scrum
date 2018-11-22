package models

// session that
type Session struct {
	AdminID     string       `json:"adminID,omitempty"`
	DateCreated string       `json:"dateCreated"`
	ID          string       `json:"ID"`
	IssueTitle  string       `json:"issueTitle"` // TODO: DELETE and move to issues
	StoryPoints []int        `json:"storyPoints"`
	Estimations []Estimation `json:"estimations,omitempty"` // TODO: DELETE and move to issues
	Issues      []Issue      `json:"issues"`
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
