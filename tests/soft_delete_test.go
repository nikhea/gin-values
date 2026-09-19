package tests

import (
	"errors"
	"net/http"
	"testing"

	"grip/config"
	"grip/dto"
	"grip/models"
	"grip/repository"
	"grip/routes"
	services "grip/service"

	"gorm.io/gorm"
)

func TestUserSoftDelete(t *testing.T) {
	requireTestDB(t)

	created, err := services.CreateUser(dto.CreateUserRequest{
		Name: "Soft User", Email: "soft@example.com", Age: 30,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if err := services.DeleteUser(created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	// Scoped reads behave as if the row is gone.
	if _, err := services.GetUserByID(created.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("get deleted err = %v, want ErrRecordNotFound", err)
	}
	if _, _, _, err := services.Login(dto.LoginRequest{Email: "soft@example.com", Password: "whatever12"}); err != services.ErrInvalidCredentials {
		t.Fatalf("login deleted err = %v, want ErrInvalidCredentials", err)
	}

	// ...but the row survives with deleted_at set.
	var row models.User
	if err := config.DB.Unscoped().First(&row, "id = ?", created.ID).Error; err != nil {
		t.Fatalf("unscoped load: %v", err)
	}
	if !row.DeletedAt.Valid {
		t.Fatal("expected deleted_at to be set")
	}

	// The freed email can be registered again (partial unique index).
	if _, err := services.CreateUser(dto.CreateUserRequest{
		Name: "Soft Again", Email: "soft@example.com", Age: 31,
	}); err != nil {
		t.Fatalf("re-register freed email: %v", err)
	}

	// Hard delete removes the row entirely.
	if err := repository.HardDeleteUser(created.ID); err != nil {
		t.Fatalf("hard delete: %v", err)
	}
	var count int64
	config.DB.Unscoped().Model(&models.User{}).Where("id = ?", created.ID).Count(&count)
	if count != 0 {
		t.Fatalf("rows = %d, want 0 after hard delete", count)
	}
}

func TestContactAndProfileSoftDelete(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()
	user, token := registerVerifiedUser(t, "Soft Owner", "softowner@example.com", "supersecret123")

	contact, err := services.CreateContact(dto.CreateContactRequest{
		UserID: user.ID, Name: "Work", Type: "email", Value: "w@example.com",
	})
	if err != nil {
		t.Fatalf("create contact: %v", err)
	}
	if err := services.DeleteContact(contact.ID); err != nil {
		t.Fatalf("delete contact: %v", err)
	}
	if _, err := services.GetContactByID(contact.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("get deleted contact err = %v, want ErrRecordNotFound", err)
	}
	var contactRow models.Contact
	if err := config.DB.Unscoped().First(&contactRow, "id = ?", contact.ID).Error; err != nil {
		t.Fatalf("unscoped contact load: %v", err)
	}
	if !contactRow.DeletedAt.Valid {
		t.Fatal("expected contact deleted_at to be set")
	}
	if err := repository.HardDeleteContact(contact.ID); err != nil {
		t.Fatalf("hard delete contact: %v", err)
	}

	// Profile soft delete through the service + HTTP flag check.
	if _, err := services.CreateProfile(user.ID, dto.CreateProfileRequest{Bio: "hi"}); err != nil {
		t.Fatalf("create profile: %v", err)
	}
	if err := services.DeleteProfile(user.ID); err != nil {
		t.Fatalf("delete profile: %v", err)
	}
	if _, err := services.GetProfileByUserID(user.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("get deleted profile err = %v, want ErrRecordNotFound", err)
	}

	w := doRequest(t, router, "DELETE", "/api/users/"+user.ID, nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("delete user status = %d", w.Code)
	}
	if got := w.Body.String(); !contains(got, `"soft_deleted":true`) {
		t.Fatalf("delete response missing soft_deleted flag: %s", got)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
