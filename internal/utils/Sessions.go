package utils

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	sessions    map[string]*domain.UserSessions
	useRedis    bool
	ctx         = context.Background()
)

// InitSessionStore initializes the session storage. If SessionStore is set to
// "redis" a Redis client is created using REDIS_ADDR, otherwise an in-memory
// map is used.
func InitSessionStore() {
	if strings.ToLower(os.Getenv("SessionStore")) == "redis" {
		useRedis = true
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}
		redisClient = redis.NewClient(&redis.Options{Addr: addr})
	} else {
		sessions = make(map[string]*domain.UserSessions)
	}
}

// SaveSession stores the session either in Redis or memory.
func SaveSession(session *domain.UserSessions) error {
	if useRedis {
		b, err := json.Marshal(session)
		if err != nil {
			return err
		}
		exp := session.Expiration
		if exp == 0 {
			exp = time.Hour
		}
		if err := redisClient.Set(ctx, "session:"+session.DisplayName, b, exp).Err(); err != nil {
			return err
		}
		if session.Jwt != "" {
			if err := redisClient.Set(ctx, "token:"+session.Jwt, session.DisplayName, exp).Err(); err != nil {
				return err
			}
		}
		return nil
	}
	sessions[session.DisplayName] = session
	return nil
}

// GetSession retrieves a session by username.
func GetSession(username string) (*domain.UserSessions, bool) {
	if useRedis {
		b, err := redisClient.Get(ctx, "session:"+username).Bytes()
		if err != nil {
			return nil, false
		}
		sess := new(domain.UserSessions)
		if err := json.Unmarshal(b, sess); err != nil {
			return nil, false
		}
		return sess, true
	}
	sess, ok := sessions[username]
	return sess, ok
}

// GetSessionByToken fetches a session using its JWT token.
func GetSessionByToken(token string) (*domain.UserSessions, bool) {
	if useRedis {
		username, err := redisClient.Get(ctx, "token:"+token).Result()
		if err != nil {
			return nil, false
		}
		return GetSession(username)
	}
	for _, v := range sessions {
		if CheckJWT(v, token) {
			return v, true
		}
	}
	return nil, false
}

// DeleteSession removes a session for the given username.
func DeleteSession(username string) {
	if useRedis {
		if sess, ok := GetSession(username); ok && sess.Jwt != "" {
			redisClient.Del(ctx, "token:"+sess.Jwt)
		}
		redisClient.Del(ctx, "session:"+username)
		return
	}
	delete(sessions, username)
}

// CheckAuthn validates the Authorization header and returns the associated session.
func CheckAuthn(c *fiber.Ctx) *domain.UserSessions {
	header := c.Get("Authorization")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil
	}
	token := parts[1]
	sess, ok := GetSessionByToken(token)
	if !ok {
		return nil
	}
	if CheckJWT(sess, token) {
		return sess
	}
	return nil
}
