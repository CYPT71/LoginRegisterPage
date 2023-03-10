package middlewares

import (
	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func CheckPerms(c *fiber.Ctx, bins uint64) error {
	session := utils.CheckAuthn(c)
	var (
		perm uint64
	)

	if session == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := domain.UserModel{}
	user.Username = session.DisplayName
	user.Get()

	perm |= (user.Permission & bins)

	
	if perm != 0 {
		return c.Next()
	} else {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
}
