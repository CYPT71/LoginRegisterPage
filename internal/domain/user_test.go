package domain

import (
	"os"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

func TestWebAuthnFields(t *testing.T) {
	u := UserModel{Id: 1, Username: "john"}
	if string(u.WebAuthnID()) == "" {
		t.Fatalf("expected id")
	}
	if u.WebAuthnName() != "john" || u.WebAuthnDisplayName() != "john" {
		t.Fatalf("wrong display or name")
	}
	if u.WebAuthnIcon() != "" {
		t.Fatalf("icon should be empty")
	}
}

func TestCredentialExcludeList(t *testing.T) {
	cred := webauthn.Credential{ID: []byte{1, 2, 3}}
	u := UserModel{Credentials: []webauthn.Credential{cred}}
	list := u.CredentialExcludeList()
	if len(list) != 1 {
		t.Fatalf("expected 1 credential")
	}
	if list[0].Type != protocol.PublicKeyCredentialType {
		t.Fatalf("wrong type")
	}
	if string(list[0].CredentialID) != string(cred.ID) {
		t.Fatalf("wrong id")
	}
}

func TestSessionExpired(t *testing.T) {

	os.Setenv("test", "true")

	// Always initialize the map
	Sessions := make(map[string]*UserSessions)

	// Create a short-lived session
	userSession := &UserSessions{
		Expiration:  time.Second * 3, // 1s for faster test
		DisplayName: "Test",
	}
	Sessions[userSession.DisplayName] = userSession

	// Start expiry timer
	go userSession.DeleteAfter(Sessions)

	// Wait a bit longer than the expiration
	time.Sleep(4 * time.Second)

	// Now the session should be gone
	if _, ok := Sessions[userSession.DisplayName]; ok {
		t.Fatal("expected session to be deleted after expiration")
	}
}
