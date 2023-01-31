package controllers

import (
	"bytes"
	"log"

	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/gofiber/fiber/v2"
)

func RegistrationStart(c *fiber.Ctx) error {

	user := new(domain.UserModel)

	user.Username = c.Params("username")

	_, ok := utils.Sessions[user.Username]

	if ok == true {
		delete(utils.Sessions, user.Username)
	}

	if user.Find() {
		log.Printf("Find user")
		return c.Status(401).JSON(fiber.Map{
			"message": "Find user",
		})
	}

	user.Permissions = 1 << 0
	user.Create()

	options, sessionData, err := utils.Web.BeginRegistration(*user)

	if err != nil {
		log.Print(err)
		return c.SendStatus(401)
	}

	session := new(domain.UserSessions)
	session.DisplayName = options.Response.User.Name
	session.SessionData = sessionData
	session.Expiration = 60 * 60 * 3600

	go session.DeleteAfter(utils.Sessions)

	utils.Sessions[session.DisplayName] = session

	return c.JSON(fiber.Map{
		"Options": options,
	})

}

func RegisterEnd(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	user.Username = c.Params("username")

	credential, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(c.Body()))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	if !user.Find() {
		return c.Status(403).JSON(fiber.Map{
			"message": "not found",
		})
	}

	session, ok := utils.Sessions[user.Username]
	log.Println(session)
	if ok == false {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "session Not exist",
		})
	}
	creds, err := utils.Web.CreateCredential(user, *session.SessionData, credential)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	user.Credentals = append(user.Credentals, *creds)

	user.SaveCredentials()

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

	return c.JSON(fiber.Map{
		"token": session.Jwt,
	})
}

func RegisterPassword(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	user.Username = c.Params("username")

	if user.Find() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"err": "not authorize",
		})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	if len(user.Password) <= 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"err": "password to short",
		})
	}

	session := new(domain.UserSessions)

	session.DisplayName = user.Username
	session.Expiration = 24 * 3600 * 2

	token, err := utils.CreateJWT(*session)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	session.Jwt = token
	utils.Sessions[user.Username] = session
	go session.DeleteAfter(utils.Sessions)

	user.Create()

	return c.JSON(fiber.Map{
		"token": session.Jwt,
	})

}
