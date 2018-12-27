package models

// Estimation represents Estimation record
type Estimation struct {
	Value   int    `json:"value"`
	IssueID string `json:"issueID"`
}

// // Validate validates the Artist fields.
// func (m Estimation) Validate() error {
// 	return validation.ValidateStruct(&m,
// 		validation.Field(&m.Username, validation.Required, validation.Length(0, 120)),
// 		validation.Field(&m.Estimate, validation.Required, is.Digit),
// 	)
// }
