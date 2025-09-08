package utils

/*
*  return true only if a session contains AAGUID
 */

import (
	"webauthn_api/internal/domain"

	"github.com/gofiber/fiber/v2"
)

var Sessions map[string]*domain.UserSessions

func CheckAuthn(c *fiber.Ctx) *domain.UserSessions {
	authType, ok := c.GetReqHeaders()["Authorization"]
	if !ok {
		return nil
	}

	if authType[0] != "Bearer" || len(authType) < 2 {
		return nil
	}

	auth := authType[1]

	// log.Println(authType)

	for _, v := range Sessions {

		if CheckJWT(v, auth) {
			return v
		}

	}
	return nil
}
