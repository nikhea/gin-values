package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"grip/dto"
	"grip/models"
	"grip/repository"
	"grip/routes"
	services "grip/service"
)

func createOrg(t *testing.T, ownerID, name string) *models.Organization {
	t.Helper()
	org, err := services.CreateOrg(ownerID, dto.CreateOrgRequest{Name: name})
	if err != nil {
		t.Fatalf("create org: %v", err)
	}
	return org
}

func TestOrgLifecycle(t *testing.T) {
	requireTestDB(t)
	owner, ownerToken := registerVerifiedUser(t, "Org Owner", "orgowner@example.com", "supersecret123")

	org := createOrg(t, owner.ID, "Acme")

	// Creator is owner in both stores.
	m, err := repository.GetMembership(owner.ID, org.ID)
	if err != nil {
		t.Fatalf("membership: %v", err)
	}
	if m.Role != models.RoleOwner {
		t.Fatalf("role = %q, want owner", m.Role)
	}

	// List shows it; a stranger sees nothing.
	orgs, err := services.ListOrgs(owner.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(orgs) != 1 || orgs[0].ID != org.ID {
		t.Fatalf("orgs = %+v", orgs)
	}
	stranger, _ := registerVerifiedUser(t, "Stranger", "stranger@example.com", "supersecret123")
	if orgs, _ := services.ListOrgs(stranger.ID); len(orgs) != 0 {
		t.Fatalf("stranger orgs = %+v", orgs)
	}

	// HTTP: owner deletes OK.
	router := routes.Setup()
	w := doRequest(t, router, "DELETE", "/api/orgs/"+org.ID, nil, ownerToken)
	if w.Code != http.StatusOK {
		t.Fatalf("owner delete status = %d, body %s", w.Code, w.Body.String())
	}
	if _, err := repository.GetOrganizationByID(org.ID); err == nil {
		t.Fatal("expected org gone")
	}
}

func TestMemberManagement(t *testing.T) {
	requireTestDB(t)
	owner, _ := registerVerifiedUser(t, "M Owner", "mowner@example.com", "supersecret123")
	member, memberToken := registerVerifiedUser(t, "M Member", "mmember@example.com", "supersecret123")
	outsider, outsiderToken := registerVerifiedUser(t, "M Outsider", "moutsider@example.com", "supersecret123")
	org := createOrg(t, owner.ID, "Globex")

	// Add member + already-member + invalid role.
	if _, err := services.AddMember(org.ID, member.ID, models.RoleMember); err != nil {
		t.Fatalf("add: %v", err)
	}
	if _, err := services.AddMember(org.ID, member.ID, models.RoleMember); err != services.ErrAlreadyMember {
		t.Fatalf("re-add err = %v, want ErrAlreadyMember", err)
	}
	if _, err := services.AddMember(org.ID, outsider.ID, "superadmin"); err != services.ErrInvalidRole {
		t.Fatalf("bad role err = %v, want ErrInvalidRole", err)
	}

	// HTTP: member cannot manage members; outsider lacks org scope (400).
	router := routes.Setup()
	w := doRequest(t, router, "POST", "/api/orgs/"+org.ID+"/members",
		dto.AddMemberRequest{UserID: outsider.ID, Role: models.RoleViewer}, memberToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("member add-member status = %d, want 403", w.Code)
	}
	w = doRequest(t, router, "DELETE", "/api/orgs/"+org.ID, nil, outsiderToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("outsider delete status = %d, want 403", w.Code)
	}

	// Two owners, then demote + remove until one remains.
	if _, err := services.UpdateMemberRole(org.ID, member.ID, models.RoleOwner); err != nil {
		t.Fatalf("promote owner: %v", err)
	}
	if _, err := services.UpdateMemberRole(org.ID, owner.ID, models.RoleMember); err != nil {
		t.Fatalf("demote with second owner present: %v", err)
	}
	if err := services.RemoveMember(org.ID, owner.ID); err != nil {
		t.Fatalf("remove demoted ex-owner: %v", err)
	}
	// Sole owner is protected both ways.
	if err := services.RemoveMember(org.ID, member.ID); err != services.ErrLastOwner {
		t.Fatalf("remove last owner err = %v, want ErrLastOwner", err)
	}
	if _, err := services.UpdateMemberRole(org.ID, member.ID, models.RoleViewer); err != services.ErrLastOwner {
		t.Fatalf("demote last owner err = %v, want ErrLastOwner", err)
	}
}

func TestContactRBACMatrix(t *testing.T) {
	requireTestDB(t)
	owner, ownerToken := registerVerifiedUser(t, "C Owner", "cowner@example.com", "supersecret123")
	member, memberToken := registerVerifiedUser(t, "C Member", "cmember@example.com", "supersecret123")
	viewer, viewerToken := registerVerifiedUser(t, "C Viewer", "cviewer@example.com", "supersecret123")
	stranger, strangerToken := registerVerifiedUser(t, "C Stranger", "cstranger@example.com", "supersecret123")
	org := createOrg(t, owner.ID, "Initech")

	for _, tc := range []struct {
		id   string
		role string
	}{
		{member.ID, models.RoleMember},
		{viewer.ID, models.RoleViewer},
	} {
		if _, err := services.AddMember(org.ID, tc.id, tc.role); err != nil {
			t.Fatalf("add %s: %v", tc.id, err)
		}
	}

	// Owner creates a team contact via HTTP (ownership forced to caller).
	router := routes.Setup()
	orgID := org.ID
	w := doRequest(t, router, "POST", "/api/contacts/",
		map[string]any{"user_id": stranger.ID, "org_id": orgID, "type": "email", "value": "team@example.com"}, ownerToken)
	if w.Code != http.StatusCreated {
		t.Fatalf("create team contact status = %d, body %s", w.Code, w.Body.String())
	}
	var created struct {
		Contact models.Contact `json:"contact"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.Contact.UserID != owner.ID {
		t.Fatalf("owner = %q, want caller %q (must not be spoofable)", created.Contact.UserID, owner.ID)
	}
	contactID := created.Contact.ID

	// Read matrix.
	for _, tc := range []struct {
		name  string
		token string
		want  int
	}{
		{"owner", ownerToken, http.StatusOK},
		{"member", memberToken, http.StatusOK},
		{"viewer", viewerToken, http.StatusOK},
		{"stranger", strangerToken, http.StatusNotFound},
	} {
		w := doRequest(t, router, "GET", "/api/contacts/"+contactID, nil, tc.token)
		if w.Code != tc.want {
			t.Fatalf("read by %s status = %d, want %d (%s)", tc.name, w.Code, tc.want, w.Body.String())
		}
	}

	// Write matrix: viewer forbidden, stranger 404, member OK.
	w = doRequest(t, router, "PUT", "/api/contacts/"+contactID,
		map[string]any{"type": "email", "value": "new@example.com"}, viewerToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("viewer write status = %d, want 403", w.Code)
	}
	w = doRequest(t, router, "PUT", "/api/contacts/"+contactID,
		map[string]any{"type": "email", "value": "new@example.com"}, strangerToken)
	if w.Code != http.StatusNotFound {
		t.Fatalf("stranger write status = %d, want 404", w.Code)
	}
	w = doRequest(t, router, "PUT", "/api/contacts/"+contactID,
		map[string]any{"type": "email", "value": "new@example.com"}, memberToken)
	if w.Code != http.StatusOK {
		t.Fatalf("member write status = %d, body %s", w.Code, w.Body.String())
	}

	// Org-scoped list: member sees it, stranger blocked.
	w = doRequest(t, router, "GET", "/api/contacts/?org_id="+orgID, nil, memberToken)
	if w.Code != http.StatusOK {
		t.Fatalf("member scoped list status = %d", w.Code)
	}
	w = doRequest(t, router, "GET", "/api/contacts/?org_id="+orgID, nil, strangerToken)
	if w.Code != http.StatusForbidden {
		t.Fatalf("stranger scoped list status = %d, want 403", w.Code)
	}

	// Unscoped list hides others' team contacts from the stranger.
	w = doRequest(t, router, "GET", "/api/contacts/", nil, strangerToken)
	if w.Code != http.StatusOK {
		t.Fatalf("stranger list status = %d", w.Code)
	}
	var list struct {
		Contacts []models.Contact `json:"contacts"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(list.Contacts) != 0 {
		t.Fatalf("stranger sees %d contacts, want 0", len(list.Contacts))
	}
}
