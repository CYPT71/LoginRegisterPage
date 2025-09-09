package utils

import (
	"os"
	"testing"
	"time"

	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestCheckAuthnMemory(t *testing.T) {
	os.Setenv("SessionStore", "memory")
	InitSessionStore()
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	session := &domain.UserSessions{DisplayName: "bob", Expiration: time.Hour}
	token, err := CreateJWT(*session)
	if err != nil {
		t.Fatalf("jwt error: %v", err)
	}
	session.Jwt = token
	SaveSession(session)
	ctx.Request().Header.Set("Authorization", "Bearer "+token)
	got := CheckAuthn(ctx)
	if got == nil || got.DisplayName != "bob" {
		t.Fatalf("invalid session")
	}
}

func TestCheckAuthnRedis(t *testing.T) {
	os.Setenv("SessionStore", "redis")
	InitSessionStore()
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	session := &domain.UserSessions{DisplayName: "alice", Expiration: time.Hour}
	token, err := CreateJWT(*session)
	if err != nil {
		t.Fatalf("jwt error: %v", err)
	}
	session.Jwt = token
	if err := SaveSession(session); err != nil {
		t.Fatalf("save error: %v", err)
	}
	ctx.Request().Header.Set("Authorization", "Bearer "+token)
	got := CheckAuthn(ctx)
	if got == nil || got.DisplayName != "alice" {
		t.Fatalf("invalid session from redis")
	}
}
