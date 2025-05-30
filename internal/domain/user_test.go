package domain

import (
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"testing"
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
