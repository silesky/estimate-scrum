package models

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Estimation represents Estimation record
type Estimation struct {
	Username string `json:"username"`
	Estimate int    `json:"estimate"`
}

// Validate validates the Artist fields.
func (m Estimation) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Username, validation.Required, validation.Length(0, 120)),
		validation.Field(&m.Estimate, validation.Required, is.Digit),
	)
}
