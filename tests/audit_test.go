package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"grip/dto"
	"grip/models"
	"grip/repository"
	"grip/routes"
	services "grip/service"

	"github.com/gin-gonic/gin"
)

func countAudit(t *testing.T, action string) int64 {
	t.Helper()
	rows, _, err := repository.ListAuditLogs(repository.AuditFilter{Action: action})
	if err != nil {
		t.Fatalf("list audit: %v", err)
	}
	return int64(len(rows))
}

func TestMutationsWriteAuditRows(t *testing.T) {
	requireTestDB(t)
	router := routes.Setup()

	// Register through HTTP (the helper bypasses handlers, where audit lives).
	w := doRequest(t, router, "POST", "/api/auth/register", dto.RegisterRequest{
		Name: "Audit User", Email: "audit@example.com", Password: "supersecret123", Age: 30,
	}, "")
	if w.Code != http.StatusCreated {
		t.Fatalf("register status = %d", w.Code)
	}
	user, err := repository.GetUserByEmail("audit@example.com")
	if err != nil {
		t.Fatalf("load user: %v", err)
	}
	if _, err := services.VerifyEmail(user.VerificationToken); err != nil {
		t.Fatalf("verify: %v", err)
	}
	w = doRequest(t, router, "POST", "/api/auth/login",
		dto.LoginRequest{Email: "audit@example.com", Password: "supersecret123"}, "")
	if w.Code != http.StatusOK {
		t.Fatalf("login status = %d", w.Code)
	}
	var loginBody struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginBody); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	token := loginBody.Token

	// Register itself wrote a row with the user as actor.
	rows, _, err := repository.ListAuditLogs(repository.AuditFilter{Action: models.AuditRegister})
	if err != nil || len(rows) != 1 {
		t.Fatalf("register rows = %d, err = %v", len(rows), err)
	}
	if rows[0].ActorID == nil || *rows[0].ActorID != user.ID {
		t.Fatalf("register actor = %v, want %s", rows[0].ActorID, user.ID)
	}
	if rows[0].RequestID == nil || *rows[0].RequestID == "" {
		t.Fatal("expected request_id captured")
	}

	// Contact create + update through HTTP.
	w = doRequest(t, router, "POST", "/api/contacts/",
		map[string]any{"type": "email", "value": "a@example.com"}, token)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status = %d", w.Code)
	}
	var created struct {
		Contact models.Contact `json:"contact"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	w = doRequest(t, router, "PUT", "/api/contacts/"+created.Contact.ID,
		map[string]any{"type": "email", "value": "b@example.com"}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("update status = %d", w.Code)
	}

	creates, _, err := repository.ListAuditLogs(repository.AuditFilter{Action: models.AuditContactCreate})
	if err != nil || len(creates) != 1 {
		t.Fatalf("create rows = %d, err = %v", len(creates), err)
	}
	if creates[0].ResourceID == nil || *creates[0].ResourceID != created.Contact.ID {
		t.Fatalf("create resource = %v, want %s", creates[0].ResourceID, created.Contact.ID)
	}
	if creates[0].ActorID == nil || *creates[0].ActorID != user.ID {
		t.Fatalf("create actor = %v, want %s", creates[0].ActorID, user.ID)
	}
	if countAudit(t, models.AuditContactUpdate) != 1 {
		t.Fatal("expected one update row")
	}

	// Import writes a summary row with counts in metadata.
	csvData := "name,type,value\nA,email,a@example.com\nB,bogus,x\n"
	w = multipartFile(t, router, "POST", "/api/contacts/import/csv", "file", "c.csv",
		[]byte(csvData), map[string]string{"user_id": user.ID}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("import status = %d", w.Code)
	}
	imports, _, err := repository.ListAuditLogs(repository.AuditFilter{Action: models.AuditContactImport})
	if err != nil || len(imports) != 1 {
		t.Fatalf("import rows = %d, err = %v", len(imports), err)
	}
	if imports[0].Metadata["imported"] != float64(1) {
		t.Fatalf("import metadata = %v", imports[0].Metadata)
	}
}

func TestAuditListingGatedByOrgAdmin(t *testing.T) {
	requireTestDB(t)
	owner, ownerToken := registerVerifiedUser(t, "A Owner", "aowner@example.com", "supersecret123")
	member, memberToken := registerVerifiedUser(t, "A Member", "amember@example.com", "supersecret123")
	outsider, outsiderToken := registerVerifiedUser(t, "A Outsider", "aoutsider@example.com", "supersecret123")
	org := createOrg(t, owner.ID, "AuditOrg")
	if _, err := services.AddMember(org.ID, member.ID, models.RoleMember); err != nil {
		t.Fatalf("add: %v", err)
	}

	router := routes.Setup()

	get := func(token, orgID string) int {
		w := doRequestOrg(t, router, "GET", "/api/audit-logs/", nil, token, orgID)
		return w.Code
	}

	// Owner (admin+) lists OK; member and outsider are forbidden.
	if code := get(ownerToken, org.ID); code != http.StatusOK {
		t.Fatalf("owner list status = %d, want 200", code)
	}
	if code := get(memberToken, org.ID); code != http.StatusForbidden {
		t.Fatalf("member list status = %d, want 403", code)
	}
	if code := get(outsiderToken, org.ID); code != http.StatusForbidden {
		t.Fatalf("outsider list status = %d, want 403", code)
	}

	// Rows are org-scoped: nothing from other orgs leaks in.
	other := createOrg(t, outsider.ID, "OtherOrg")
	w := doRequestOrg(t, router, "POST", "/api/orgs/"+other.ID+"/members",
		dto.AddMemberRequest{UserID: member.ID, Role: models.RoleViewer}, outsiderToken, other.ID)
	if w.Code != http.StatusCreated {
		t.Fatalf("seed other-org action status = %d", w.Code)
	}
	rows, _, err := repository.ListAuditLogs(repository.AuditFilter{OrgID: org.ID})
	if err != nil {
		t.Fatalf("scoped list: %v", err)
	}
	for _, r := range rows {
		if r.OrgID != nil && *r.OrgID != org.ID {
			t.Fatalf("leaked row from org %s", *r.OrgID)
		}
	}

	// Detail endpoint respects the same scope.
	all, _, err := repository.ListAuditLogs(repository.AuditFilter{OrgID: other.ID})
	if err != nil || len(all) == 0 {
		t.Fatalf("need at least one other-org row, got %d (%v)", len(all), err)
	}
	w = doRequestOrg(t, router, "GET", "/api/audit-logs/"+all[0].ID, nil, ownerToken, org.ID)
	if w.Code != http.StatusNotFound {
		t.Fatalf("cross-org detail status = %d, want 404", w.Code)
	}
	w = doRequestOrg(t, router, "GET", "/api/audit-logs/"+all[0].ID, nil, outsiderToken, other.ID)
	if w.Code != http.StatusOK {
		t.Fatalf("own-org detail status = %d, want 200", w.Code)
	}
}

// doRequestOrg is doRequest with an X-Org-ID scope header.
func doRequestOrg(t *testing.T, router *gin.Engine, method, path string, body any, token, orgID string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if orgID != "" {
		req.Header.Set("X-Org-ID", orgID)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
