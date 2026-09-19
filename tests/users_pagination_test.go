package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"grip/dto"
	"grip/models"
	"grip/routes"
	services "grip/service"
)

func seedUsers(t *testing.T, users ...dto.CreateUserRequest) {
	t.Helper()
	for _, u := range users {
		if _, err := services.CreateUser(u); err != nil {
			t.Fatalf("seed %s: %v", u.Email, err)
		}
	}
}

func TestListUsersSearch(t *testing.T) {
	requireTestDB(t)
	seedUsers(t,
		dto.CreateUserRequest{Name: "Alpha One", Email: "alpha@example.com", Age: 30},
		dto.CreateUserRequest{Name: "Beta Two", Email: "beta@example.com", Age: 30},
		dto.CreateUserRequest{Name: "Alpha Three", Email: "alpha3@example.com", Age: 30},
	)

	users, total, err := services.ListUsers(dto.UserFilter{Search: "alpha"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 2 || len(users) != 2 {
		t.Fatalf("total = %d rows = %d, want 2/2", total, len(users))
	}

	users, total, err = services.ListUsers(dto.UserFilter{Search: "BETA@EXAMPLE"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(users) != 1 {
		t.Fatalf("case-insensitive total = %d, want 1", total)
	}

	users, total, err = services.ListUsers(dto.UserFilter{Search: "zzz-no-match"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 0 || len(users) != 0 {
		t.Fatalf("total = %d, want 0", total)
	}
}

func TestListUsersPages(t *testing.T) {
	requireTestDB(t)
	seedUsers(t,
		dto.CreateUserRequest{Name: "Page One", Email: "p1@example.com", Age: 30},
		dto.CreateUserRequest{Name: "Page Two", Email: "p2@example.com", Age: 30},
		dto.CreateUserRequest{Name: "Page Three", Email: "p3@example.com", Age: 30},
	)

	users, total, err := services.ListUsers(dto.UserFilter{Pagination: dto.Pagination{Page: 1, PageSize: 2}})
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	if total != 3 || len(users) != 2 {
		t.Fatalf("page 1: total = %d rows = %d, want 3/2", total, len(users))
	}

	users, total, err = services.ListUsers(dto.UserFilter{Pagination: dto.Pagination{Page: 2, PageSize: 2}})
	if err != nil {
		t.Fatalf("page 2: %v", err)
	}
	if total != 3 || len(users) != 1 {
		t.Fatalf("page 2: total = %d rows = %d, want 3/1", total, len(users))
	}

	users, total, err = services.ListUsers(dto.UserFilter{Pagination: dto.Pagination{Page: 9, PageSize: 10}})
	if err != nil {
		t.Fatalf("page 9: %v", err)
	}
	if total != 3 || len(users) != 0 {
		t.Fatalf("page 9: total = %d rows = %d, want 3/0", total, len(users))
	}
}

func TestListUsersHTTP(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()
	_, token := registerVerifiedUser(t, "Page HTTP", "pagehttp@example.com", "supersecret123")
	seedUsers(t,
		dto.CreateUserRequest{Name: "HTTP Alpha", Email: "httpalpha@example.com", Age: 30},
		dto.CreateUserRequest{Name: "HTTP Beta", Email: "httpbeta@example.com", Age: 30},
	)

	w := doRequest(t, router, "GET", "/api/users/?search=alpha&page=1&page_size=1", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", w.Code, w.Body.String())
	}
	var body struct {
		Users []models.User `json:"user"`
		Meta  dto.PageMeta  `json:"meta"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Users) != 1 || body.Meta.Total < 1 || body.Meta.Page != 1 || body.Meta.PageSize != 1 {
		t.Fatalf("body = %+v", body)
	}
	if body.Meta.TotalPages != int((body.Meta.Total+1-1)/1) {
		t.Fatalf("bad total pages: %+v", body.Meta)
	}

	w = doRequest(t, router, "GET", "/api/users/?page=bogus", nil, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("bad page status = %d, want 400", w.Code)
	}
}
