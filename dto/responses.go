package dto

import "gin-learn/models"

// Response envelopes mirror the JSON shapes returned by handlers so
// Swagger documents accurate schemas instead of generic objects.

type UserEnvelope struct {
	Message string      `json:"message"`
	User    models.User `json:"user"`
}

type UsersEnvelope struct {
	Message string        `json:"message"`
	Users   []models.User `json:"user"`
}

type ProfileEnvelope struct {
	Message string         `json:"message"`
	Profile models.Profile `json:"profile"`
}

type ContactEnvelope struct {
	Message string         `json:"message"`
	Contact models.Contact `json:"contact"`
}

type ContactsEnvelope struct {
	Message  string           `json:"message"`
	Contacts []models.Contact `json:"contacts"`
	Meta     PageMeta         `json:"meta"`
}

type MessageEnvelope struct {
	Message string `json:"message"`
}

type ErrorEnvelope struct {
	Error string `json:"error"`
}
