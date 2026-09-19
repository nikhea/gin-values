package dto

// UserFilter carries the users list endpoint's filter query params.
type UserFilter struct {
	Pagination
	// Search matches name and email (case-insensitive substring).
	Search string `form:"search"`
}
