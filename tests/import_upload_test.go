package tests

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"grip/dto"
	"grip/repository"
	"grip/routes"
	services "grip/service"

	"github.com/gin-gonic/gin"
)

// tinyPNG is a valid 1x1 PNG used for avatar upload tests.
var tinyPNG = mustHex("89504e470d0a1a0a0000000d49484452000000010000000108060000001f15c4890000000d49444154789c626001000000ffff03000006000557bfabd40000000049454e44ae426082")

func mustHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

// multipartFile builds a multipart/form-data request with one file field
// plus extra form fields.
func multipartFile(t *testing.T, router *gin.Engine, method, path, fileField, filename string, content []byte, fields map[string]string, token string) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			t.Fatalf("write field: %v", err)
		}
	}
	fw, err := w.CreateFormFile(fileField, filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(method, path, &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func decodeSummary(t *testing.T, rec *httptest.ResponseRecorder) dto.ImportSummary {
	t.Helper()
	var summary dto.ImportSummary
	if err := json.Unmarshal(rec.Body.Bytes(), &summary); err != nil {
		t.Fatalf("decode summary: %v (body %s)", err, rec.Body.String())
	}
	return summary
}

func TestContactCSVImport(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()
	user, token := registerVerifiedUser(t, "Import User", "csvimport@example.com", "supersecret123")

	csvData := "name,type,value\nWork email,email,kaige@work.com\nMobile,phone,+15551234567\nBad,bogus,xyz\n"
	w := multipartFile(t, router, "POST", "/api/contacts/import/csv", "file", "contacts.csv",
		[]byte(csvData), map[string]string{"user_id": user.ID}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("csv import status = %d, body %s", w.Code, w.Body.String())
	}
	summary := decodeSummary(t, w)
	if summary.Imported != 2 || summary.Failed != 1 {
		t.Fatalf("summary = %+v, want imported=2 failed=1", summary)
	}
	if len(summary.Errors) != 1 || summary.Errors[0].Row != 4 {
		t.Fatalf("errors = %+v, want one error at row 4", summary.Errors)
	}

	contacts, total, err := services.ListContacts(dto.ContactFilter{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 2 || len(contacts) != 2 {
		t.Fatalf("total = %d, want 2", total)
	}

	// No auth token.
	w = multipartFile(t, router, "POST", "/api/contacts/import/csv", "file", "contacts.csv",
		[]byte(csvData), map[string]string{"user_id": user.ID}, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("no-token status = %d, want 401", w.Code)
	}

	// Bad header.
	w = multipartFile(t, router, "POST", "/api/contacts/import/csv", "file", "contacts.csv",
		[]byte("a,b\n1,2\n"), map[string]string{"user_id": user.ID}, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("bad header status = %d, want 400", w.Code)
	}

	// Wrong extension.
	w = multipartFile(t, router, "POST", "/api/contacts/import/csv", "file", "contacts.txt",
		[]byte(csvData), map[string]string{"user_id": user.ID}, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("bad ext status = %d, want 400", w.Code)
	}

	// Missing user_id.
	w = multipartFile(t, router, "POST", "/api/contacts/import/csv", "file", "contacts.csv",
		[]byte(csvData), nil, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("missing user status = %d, want 400", w.Code)
	}

	// Unknown user.
	w = multipartFile(t, router, "POST", "/api/contacts/import/csv", "file", "contacts.csv",
		[]byte(csvData), map[string]string{"user_id": "no-such-user"}, token)
	if w.Code != http.StatusNotFound {
		t.Fatalf("unknown user status = %d, want 404", w.Code)
	}
}

func TestContactJSONImport(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()
	user, token := registerVerifiedUser(t, "JSON Import", "jsonimport@example.com", "supersecret123")

	jsonData := `[{"name":"Work email","type":"email","value":"kaige@work.com"},{"name":"Bad","type":"fax","value":"123"}]`
	w := multipartFile(t, router, "POST", "/api/contacts/import/json", "file", "contacts.json",
		[]byte(jsonData), map[string]string{"user_id": user.ID}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("json import status = %d, body %s", w.Code, w.Body.String())
	}
	summary := decodeSummary(t, w)
	if summary.Imported != 1 || summary.Failed != 1 {
		t.Fatalf("summary = %+v, want imported=1 failed=1", summary)
	}

	// Malformed JSON.
	w = multipartFile(t, router, "POST", "/api/contacts/import/json", "file", "contacts.json",
		[]byte(`[{"name":`), map[string]string{"user_id": user.ID}, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("malformed status = %d, want 400", w.Code)
	}

	// Empty array.
	w = multipartFile(t, router, "POST", "/api/contacts/import/json", "file", "contacts.json",
		[]byte(`[]`), map[string]string{"user_id": user.ID}, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("empty array status = %d, want 400", w.Code)
	}
}

func TestUserAvatarUpload(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()
	user, token := registerVerifiedUser(t, "Avatar User", "avatar@example.com", "supersecret123")

	// No profile yet: upload creates one.
	w := multipartFile(t, router, "POST", "/api/users/"+user.ID+"/avatar", "avatar", "me.png", tinyPNG, nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("avatar status = %d, body %s", w.Code, w.Body.String())
	}
	var resp dto.AvatarResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil || resp.AvatarURL == "" {
		t.Fatalf("bad avatar response: %s", w.Body.String())
	}
	if !strings.HasPrefix(resp.AvatarURL, "/uploads/avatars/users/") || !strings.HasSuffix(resp.AvatarURL, ".png") {
		t.Fatalf("avatar_url = %q", resp.AvatarURL)
	}
	if _, err := os.Stat(filepath.Join(".", resp.AvatarURL)); err != nil {
		t.Fatalf("avatar file missing on disk: %v", err)
	}
	profile, err := repository.GetProfileByUserID(user.ID)
	if err != nil {
		t.Fatalf("load profile: %v", err)
	}
	if profile.AvatarURL != resp.AvatarURL {
		t.Fatalf("profile avatar = %q, want %q", profile.AvatarURL, resp.AvatarURL)
	}

	// Static serving exposes the file.
	req := httptest.NewRequest(http.MethodGet, resp.AvatarURL, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET avatar status = %d, want 200", rec.Code)
	}

	// Second upload with existing profile still works.
	w = multipartFile(t, router, "POST", "/api/users/"+user.ID+"/avatar", "avatar", "me2.jpg", tinyPNG, nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("second avatar status = %d, body %s", w.Code, w.Body.String())
	}

	// Non-image rejected.
	w = multipartFile(t, router, "POST", "/api/users/"+user.ID+"/avatar", "avatar", "evil.txt", []byte("hi"), nil, token)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("txt avatar status = %d, want 400", w.Code)
	}

	// Unknown user.
	w = multipartFile(t, router, "POST", "/api/users/no-such-user/avatar", "avatar", "me.png", tinyPNG, nil, token)
	if w.Code != http.StatusNotFound {
		t.Fatalf("unknown user avatar status = %d, want 404", w.Code)
	}
}

func TestContactAvatarUpload(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()
	user, token := registerVerifiedUser(t, "Contact Avatar", "cavatar@example.com", "supersecret123")

	contact, err := services.CreateContact(dto.CreateContactRequest{
		UserID: user.ID, Name: "Work", Type: "email", Value: "kaige@work.com",
	})
	if err != nil {
		t.Fatalf("create contact: %v", err)
	}

	w := multipartFile(t, router, "POST", "/api/contacts/"+contact.ID+"/avatar", "avatar", "c.png", tinyPNG, nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("avatar status = %d, body %s", w.Code, w.Body.String())
	}
	var resp dto.AvatarResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !strings.HasPrefix(resp.AvatarURL, "/uploads/avatars/contacts/") {
		t.Fatalf("avatar_url = %q", resp.AvatarURL)
	}

	stored, err := repository.GetContactByID(contact.ID)
	if err != nil {
		t.Fatalf("load contact: %v", err)
	}
	if stored.AvatarURL != resp.AvatarURL {
		t.Fatalf("contact avatar = %q, want %q", stored.AvatarURL, resp.AvatarURL)
	}

	// Unknown contact.
	w = multipartFile(t, router, "POST", "/api/contacts/no-such-id/avatar", "avatar", "c.png", tinyPNG, nil, token)
	if w.Code != http.StatusNotFound {
		t.Fatalf("unknown contact status = %d, want 404", w.Code)
	}
}
