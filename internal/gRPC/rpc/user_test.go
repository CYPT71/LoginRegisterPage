package rpc

import (
	"context"
	"testing"

	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"google.golang.org/grpc/metadata"
)

func TestSessionFromContext(t *testing.T) {
	utils.InitSessionStore()
	sess := &domain.UserSessions{DisplayName: "bob"}
	token, err := utils.CreateJWT(*sess)
	if err != nil {
		t.Fatalf("CreateJWT failed: %v", err)
	}
	sess.Jwt = token
	utils.SaveSession(sess)
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	got, err := sessionFromContext(ctx)
	if err != nil {
		t.Fatalf("sessionFromContext returned error: %v", err)
	}
	if got.DisplayName != "bob" {
		t.Fatalf("expected bob, got %s", got.DisplayName)
	}
}

func TestSessionFromContextMissing(t *testing.T) {
	ctx := context.Background()
	if _, err := sessionFromContext(ctx); err == nil {
		t.Fatalf("expected error")
	}
}
