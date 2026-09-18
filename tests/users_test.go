package tests

import (
	"errors"
	"testing"

	"grip/dto"
	services "grip/service"

	"gorm.io/gorm"
)

func TestUserCRUDService(t *testing.T) {
	requireTestDB(t)

	created, err := services.CreateUser(dto.CreateUserRequest{
		Name: "Crud User", Email: "crud@example.com", Age: 33,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected ID to be assigned")
	}

	users, err := services.GetUsers()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("got %d users, want 1", len(users))
	}

	got, err := services.GetUserByID(created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Email != "crud@example.com" {
		t.Fatalf("email = %q", got.Email)
	}

	updated, err := services.UpdateUser(created.ID, dto.CreateUserRequest{
		Name: "Crud Renamed", Email: "crud@example.com", Age: 34,
	})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Crud Renamed" || updated.Age != 34 {
		t.Fatalf("unexpected updated user: %+v", updated)
	}

	if err := services.DeleteUser(created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := services.GetUserByID(created.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("get deleted err = %v, want ErrRecordNotFound", err)
	}
	if err := services.DeleteUser("missing-id"); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("delete missing err = %v, want ErrRecordNotFound", err)
	}
}
