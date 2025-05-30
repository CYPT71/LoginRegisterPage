package utils

import (
	"github.com/duo-labs/webauthn/webauthn"
	"testing"
	"webauthn_api/internal/domain"
)

func TestCheckJWTValid(t *testing.T) {
	session := domain.UserSessions{
		DisplayName: "test",
		SessionData: &webauthn.SessionData{UserID: []byte("id")},
	}
	token, err := CreateJWT(session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !CheckJWT(&session, token) {
		t.Fatalf("expected token to be valid")
	}
}

func TestCheckJWTNilSessionData(t *testing.T) {
	session := domain.UserSessions{DisplayName: "test"}
	token, err := CreateJWT(session)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !CheckJWT(&session, token) {
		t.Fatalf("expected token to be valid even without session data")
	}
}

func TestCreateJWT(t *testing.T) {
	session := domain.UserSessions{DisplayName: "alice"}
	token, err := CreateJWT(session)
	if err != nil || token == "" {
		t.Fatalf("token should be created")
	}
}
