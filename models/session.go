package models

// Session is Session record
type Session struct {
	DateCreated string
	ID          string
	AdminID     string
	IssueTitle  string // optional issue
	StoryPoints []int  // represents story point values for a given session
}

// database contains a bunch of these sessions, organized like this
/*
id1: {...mySession},
id2: {...myOtherSession},

*/
