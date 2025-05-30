package utils

import (
	"github.com/gofiber/fiber/v2"
	"testing"
	"webauthn_api/internal/domain"
)

func TestCheckAuthn(t *testing.T) {
	app := fiber.New()
	req := app.AcquireCtx(&fiber.Ctx{})
	session := &domain.UserSessions{DisplayName: "bob"}
	token, err := CreateJWT(*session)
	if err != nil {
		t.Fatalf("jwt error: %v", err)
	}
	session.Jwt = token
	Sessions = map[string]*domain.UserSessions{"bob": session}
	req.Request().Header.Set("Authorization", "Bearer "+token)
	got := CheckAuthn(req)
	if got == nil || got.DisplayName != "bob" {
		t.Fatalf("invalid session")
	}
}
