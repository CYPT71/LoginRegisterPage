package main

import (
	"encoding/json"
	"strings"

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
	user := new(UserModel)
	user.Username = c.Params("username")

	return c.Status(200).JSON(fiber.Map{
		"user": user.Get() != nil,
	})

}

func UserBootstrap(app fiber.Router) {

	app.Get("/", func(c *fiber.Ctx) error {
		user := new(UserModel)
		userSession := checkAuthn(c)
		user.Username = userSession.displayName
		return c.Status(200).JSON(user.Get())
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		userSession := checkAuthn(c)
		delete(sessions, userSession.displayName)
		return c.Status(200).JSON(fiber.Map{
			"message": "logout",
		})
	})

	app.Patch("/", func(c *fiber.Ctx) error {

		userSession := checkAuthn(c)

		userIn := new(UserModel)
		userIn.Username = userSession.displayName
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
		user := new(UserModel)
		userSession := checkAuthn(c)
		user.Username = userSession.displayName

		user.Delete()
		delete(sessions, user.Username)

		return c.JSON(fiber.Map{
			"message": "deleted",
		})
	})

	app.Delete("/cred", func(c *fiber.Ctx) error {
		user := new(UserModel)
		userSession := checkAuthn(c)
		user.Username = userSession.displayName
		user = user.Get()
		if user == nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		user.Credentials = strings.Split(user.Credentials, ";")[0]
		user.Update()

		return c.Status(200).JSON(user)
	})

}
