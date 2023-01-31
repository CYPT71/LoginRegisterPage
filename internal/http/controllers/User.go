package controllers

import (
	"encoding/json"
	"strings"

	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type PartialUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *PartialUser) Unmarshal(body []byte) error {
	return json.Unmarshal(body, &p)
}

func CheckUserName(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	user.Username = c.Params("username")

	return c.Status(200).JSON(fiber.Map{
		"user": user.Get() != nil,
	})

}

func UserBootstrap(app fiber.Router) {

	app.Get("/", func(c *fiber.Ctx) error {
		user := new(domain.UserModel)
		userSession := utils.CheckAuthn(c)
		user.Username = userSession.DisplayName
		return c.Status(200).JSON(user.Get())
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		userSession := utils.CheckAuthn(c)
		delete(utils.Sessions, userSession.DisplayName)
		return c.Status(200).JSON(fiber.Map{
			"message": "logout",
		})
	})

	app.Patch("/", func(c *fiber.Ctx) error {

		userSession := utils.CheckAuthn(c)

		userIn := new(domain.UserModel)
		userIn.Username = userSession.DisplayName
		userIn = userIn.Get()

		
		user := new(PartialUser)
		err := user.Unmarshal(c.Body())
		if err != nil {
			return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"err": err.Error(),
			})
		}

		userIn.Email = user.Email
		userIn.Password = user.Password

		userIn.Update()

		return c.Status(200).JSON(user)

	})

	app.Delete("/", func(c *fiber.Ctx) error {
		user := new(domain.UserModel)
		userSession := utils.CheckAuthn(c)
		user.Username = userSession.DisplayName

		user.Delete()
		delete(utils.Sessions, user.Username)

		return c.JSON(fiber.Map{
			"message": "deleted",
		})
	})

	app.Delete("/cred", func(c *fiber.Ctx) error {
		user := new(domain.UserModel)
		userSession := utils.CheckAuthn(c)
		user.Username = userSession.DisplayName
		user = user.Get()
		if user == nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		user.Credentials = strings.Split(user.Credentials, ";")[0]
		user.Update()

		return c.Status(200).JSON(user)
	})

}
