package http

import (
	"log"
	"webauthn_api/internal/http/controllers"
	"webauthn_api/internal/utils"

	"webauthn_api/internal/http/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Http() *fiber.App {
	app := fiber.New()

	app.Use(middlewares.CORS())

	// app.Get("/checkUser/:username", CheckUserName)
	//app routes
	app.Post("register/start/:username", controllers.RegistrationStart)

	app.Post("register/end/:username", controllers.RegisterEnd)

	app.Post("register/password/:username", controllers.RegisterPassword)

	app.Post("login/start/:username", controllers.LoginStart)

	app.Post("login/end/:username", controllers.LoginEnd)

	app.Post("login/password/:username", controllers.LoginPassword)

	controllers.UserBootstrap(app.Group("user", func(c *fiber.Ctx) error {

		if utils.CheckAuthn(c) == nil {
			log.Println(c.GetReqHeaders())
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	}))

	return app
}
