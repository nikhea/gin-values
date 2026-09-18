package dto

import (
	"fmt"
	"net/mail"
	"regexp"
)

var phoneRegexp = regexp.MustCompile(`^\+?[0-9][0-9\s\-().]{5,20}$`)

// CreateContactRequest creates a contact for a user.
type CreateContactRequest struct {
	UserID string `json:"user_id" binding:"required,max=36" example:"458622d8-daba-4252-8ce1-846277353139"`

	Name string `json:"name" binding:"omitempty,max=100" example:"Work email"`

	Type string `json:"type" binding:"required,oneof=email phone address other" example:"email" enums:"email,phone,address,other"`

	Value string `json:"value" binding:"required,max=2048" example:"kaige@work.com"`
}

// UpdateContactRequest replaces a contact's editable fields.
// UserID is immutable and therefore not part of the payload.
type UpdateContactRequest struct {
	Name string `json:"name" binding:"omitempty,max=100" example:"Work email"`

	Type string `json:"type" binding:"required,oneof=email phone address other" example:"phone" enums:"email,phone,address,other"`

	Value string `json:"value" binding:"required,max=2048" example:"+15551234567"`
}

// ContactFilter carries the list endpoint's filter query params.
type ContactFilter struct {
	Pagination
	UserID string `form:"user_id"`
	Type   string `form:"type" binding:"omitempty,oneof=email phone address other"`
	Search string `form:"search"`
}

// ValidationError marks client-caused validation failures so handlers
// can answer 400 instead of 500. Use errors.As to detect it.
type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string { return e.Msg }

// ValidateContactValue checks Value against its declared Type:
// email must parse, phone must look like an international number,
// address/other just need a non-empty value.
func ValidateContactValue(contactType, value string) error {
	switch contactType {
	case "email":
		if _, err := mail.ParseAddress(value); err != nil {
			return &ValidationError{Msg: fmt.Sprintf("invalid email value: %q", value)}
		}
	case "phone":
		if !phoneRegexp.MatchString(value) {
			return &ValidationError{Msg: fmt.Sprintf("invalid phone value: %q (expected digits, e.g. +15551234567)", value)}
		}
	case "address", "other":
		if value == "" {
			return &ValidationError{Msg: fmt.Sprintf("value must not be empty for type %q", contactType)}
		}
	default:
		return &ValidationError{Msg: fmt.Sprintf("unknown contact type: %q", contactType)}
	}
	return nil
}
