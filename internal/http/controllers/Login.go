package controllers

import (
	"bytes"
	"log"

	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/gofiber/fiber/v2"
)

func LoginStart(c *fiber.Ctx) error {

	user := new(domain.UserModel)
	user.Username = c.Params("username")
	if user.Find() == false {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "no user with this username",
		})
	}
	user.ParseCredentials()

	options, sessionData, err := utils.Web.BeginLogin(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	session := new(domain.UserSessions)

	session.SessionData = sessionData
	session.DisplayName = user.Username
	session.Expiration = 60 * 3
	go session.DeleteAfter(utils.Sessions)
	utils.Sessions[user.Username] = session

	return c.JSON(options)

}

func LoginEnd(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	user.Username = c.Params("username")
	if user.Find() == false {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "no user with this username",
		})
	}
	user.ParseCredentials()

	session, ok := utils.Sessions[c.Params("username")]

	if ok == false {

		log.Println("session Not exist")
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "session Not exist",
		})
	}

	parsedCredential, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	creds, err := utils.Web.ValidateLogin(user, *session.SessionData, parsedCredential)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}
	session.SessionCred = creds
	session.Expiration = 24 * 3600 * 2
	token, err := utils.CreateJWT(*session)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	session.Jwt = token
	go session.DeleteAfter(utils.Sessions)

	utils.Sessions[user.Username] = session
	user.Credentals = append(user.Credentals, *creds)

	user.SaveCredentials()

	return c.JSON(fiber.Map{
		"token": session.Jwt,
	})
}

func LoginPassword(c *fiber.Ctx) error {

	user := new(domain.UserModel)
	user.Username = c.Params("username")
	if user.Find() == false {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "no user with this username",
		})
	}

	userBody := new(domain.UserModel)

	if err := c.BodyParser(&userBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	if userBody.Password != user.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": "Not Authorize",
		})
	}

	session := new(domain.UserSessions)

	session.DisplayName = user.Username
	token, err := utils.CreateJWT(*session)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	session.Jwt = token
	session.Expiration = 24 * 3600 * 2
	go session.DeleteAfter(utils.Sessions)

	utils.Sessions[user.Username] = session

	return c.JSON(fiber.Map{
		"token": session.Jwt,
	})
}
