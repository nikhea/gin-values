package tests

import (
	"errors"
	"testing"

	"gin-learn/dto"
)

func TestValidateContactValue(t *testing.T) {
	tests := []struct {
		name        string
		contactType string
		value       string
		wantErr     bool
	}{
		{"valid email", "email", "user@example.com", false},
		{"invalid email", "email", "not-an-email", true},
		{"valid phone", "phone", "+15551234567", false},
		{"invalid phone", "phone", "abc", true},
		{"valid address", "address", "123 Main St", false},
		{"empty address", "address", "", true},
		{"valid other", "other", "anything", false},
		{"empty other", "other", "", true},
		{"unknown type", "fax", "123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := dto.ValidateContactValue(tt.contactType, tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateContactValue(%q, %q) err = %v, wantErr = %v",
					tt.contactType, tt.value, err, tt.wantErr)
			}
			if err != nil {
				var vErr *dto.ValidationError
				if !errors.As(err, &vErr) {
					t.Fatalf("expected *ValidationError, got %T", err)
				}
			}
		})
	}
}
