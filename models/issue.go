package models

// e.g. {[userName]: point total}
type EstimationsMap = map[string]int
type Issue struct {
	IssueTitle  string         `json:"issueTitle"`
	IssueID     string         `json:"issueID"`
	Estimations EstimationsMap `json:"estimations"` // userID: 123
}
