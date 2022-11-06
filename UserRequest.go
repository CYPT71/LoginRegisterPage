package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

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
		userIn := new(UserModel)

		user := new(UserModel)
		if err := c.BodyParser(user); err != nil {
			fmt.Println("error = ", err)
			return c.SendStatus(200)
		}
		userSession := checkAuthn(c)

		userIn.Username = userSession.displayName
		userIn = userIn.Get()

		user.Username = userSession.displayName
		user.Credentials = userIn.Credentials
		// user.Update()

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
