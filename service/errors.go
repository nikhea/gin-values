package services

import "errors"

// ErrForbidden marks authorization failures; handlers map it to 403.
// Unknown/non-member resources should usually surface as 404 instead
// (gorm.ErrRecordNotFound) to avoid leaking existence.
var ErrForbidden = errors.New("forbidden")
