package models

import validation "github.com/go-ozzo/ozzo-validation"

// Estimation represents Estimation record
type Estimation struct {
	Value   int    `json:"value"`
	IssueID string `json:"issueID"`
}

func (s Estimation) OK() error {
	return validation.ValidateStruct(
		validation.Field(s.IssueID, validation.Required),
		validation.Field(s.Value, validation.Required),
		validation.Field(s.Value, validation.Max(15)),
	)
}

// // Validate validates the Artist fields.
// func (m Estimation) Validate() error {
// 	return validation.ValidateStruct(&m,
// 		validation.Field(&m.Username, validation.Required, validation.Length(0, 120)),
// 		validation.Field(&m.Estimate, validation.Required, is.Digit),
// 	)
// }
