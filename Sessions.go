package main

import (
	"strings"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
)

type UserSessions struct {
	sessionData *webauthn.SessionData `json:"-"`
	sessionCred *webauthn.Credential  `json:"-"`
	displayName string
	jwt         string
	expiration  uint64 `json:"-"`
}

/*
*  return true only if a session contains AAGUID
 */
func checkAuthn(c *fiber.Ctx) *UserSessions {
	value, ok := c.GetReqHeaders()["Authorization"]
	if ok == false {
		return nil
	}
	authType := strings.Split(value, " ")
	if authType[0] != "Bearer" || len(authType) < 2 {
		return nil
	}

	auth := authType[1]

	for _, v := range sessions {

		if checkJWT(v, auth) {
			return v
		}

	}
	return nil
}
