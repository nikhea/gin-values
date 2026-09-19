package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"grip/dto"
	"grip/jobs"
	"grip/models"
	"grip/repository"
	"grip/routes"
)

// waitForNotification polls until the worker writes the row (River works
// async) or the timeout expires.
func waitForNotification(t *testing.T, userID, typ string) models.Notification {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		rows, _, err := repository.ListNotifications(repository.NotificationFilter{UserID: userID})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		for _, n := range rows {
			if n.Type == typ {
				return n
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %q notification", typ)
	return models.Notification{}
}

func TestNotifyWorkerWritesRow(t *testing.T) {
	requireTestDB(t)
	user, _ := registerVerifiedUser(t, "Notify Worker", "notifyworker@example.com", "supersecret123")

	if err := jobs.EnqueueNotify(context.Background(), user.ID, models.NotifyWelcome,
		"Hello", "World", map[string]any{"k": "v"}); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	n := waitForNotification(t, user.ID, models.NotifyWelcome)
	if n.Title != "Hello" || n.Body != "World" {
		t.Fatalf("row = %+v", n)
	}
	if n.Data["k"] != "v" {
		t.Fatalf("data = %v", n.Data)
	}
	if n.ReadAt != nil {
		t.Fatal("expected unread")
	}
}

func TestInboxFlow(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()
	user, token := registerVerifiedUser(t, "Inbox User", "inbox@example.com", "supersecret123")
	_, otherToken := registerVerifiedUser(t, "Inbox Other", "inboxother@example.com", "supersecret123")

	// Welcome notification arrives via the verify flow (HTTP register).
	w := doRequest(t, router, "POST", "/api/auth/register", dto.RegisterRequest{
		Name: "Inbox HTTP", Email: "inboxhttp@example.com", Password: "supersecret123", Age: 30,
	}, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("register status = %d", w.Code)
	}
	httpUser, err := repository.GetUserByEmail("inboxhttp@example.com")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	w = doRequest(t, router, "GET", "/api/auth/verify?token="+httpUser.VerificationToken, nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("verify status = %d", w.Code)
	}

	// The HTTP verify handler emits the welcome through River.
	waitForNotification(t, httpUser.ID, models.NotifyWelcome)

	// Inbox flow for the helper user, seeded directly.
	if err := jobs.EnqueueNotify(context.Background(), user.ID, models.NotifyWelcome,
		"Hello", "World", nil); err != nil {
		t.Fatalf("enqueue: %v", err)
	}
	waitForNotification(t, user.ID, models.NotifyWelcome)

	// List + unread count.
	w = doRequest(t, router, "GET", "/api/notifications/", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("list status = %d, body %s", w.Code, w.Body.String())
	}
	var list struct {
		Notifications []models.Notification `json:"notifications"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(list.Notifications) == 0 {
		t.Fatal("expected inbox rows")
	}
	w = doRequest(t, router, "GET", "/api/notifications/unread-count", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("count status = %d", w.Code)
	}
	var count struct {
		Unread int64 `json:"unread"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &count); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if count.Unread != int64(len(list.Notifications)) {
		t.Fatalf("unread = %d, want %d", count.Unread, len(list.Notifications))
	}

	// Mark one read.
	first := list.Notifications[0].ID
	w = doRequest(t, router, "PATCH", "/api/notifications/"+first+"/read", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("mark read status = %d", w.Code)
	}

	// Other user can't read it (404, no leak) and has an empty inbox.
	w = doRequest(t, router, "PATCH", "/api/notifications/"+first+"/read", nil, otherToken)
	if w.Code != http.StatusNotFound {
		t.Fatalf("cross-user read status = %d, want 404", w.Code)
	}
	w = doRequest(t, router, "GET", "/api/notifications/", nil, otherToken)
	var otherList struct {
		Notifications []models.Notification `json:"notifications"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &otherList); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(otherList.Notifications) != 0 {
		t.Fatalf("other inbox = %d rows, want 0", len(otherList.Notifications))
	}

	// Read-all is idempotent.
	w = doRequest(t, router, "POST", "/api/notifications/read-all", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("read-all status = %d", w.Code)
	}
	w = doRequest(t, router, "POST", "/api/notifications/read-all", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("read-all again status = %d", w.Code)
	}
	w = doRequest(t, router, "GET", "/api/notifications/unread-count", nil, token)
	var after struct {
		Unread int64 `json:"unread"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &after); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if after.Unread != 0 {
		t.Fatalf("unread after read-all = %d, want 0", after.Unread)
	}
}
