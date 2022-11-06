package main

import (
	"bytes"
	"log"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/gofiber/fiber/v2"
)

func (session *UserSessions) deleteAfter() {
	for i := session.expiration; i >= 0; i-- {
		time.Sleep(1)
	}

	log.Printf("user delete")
	delete(sessions, session.displayName)
}

func RegistrationStart(c *fiber.Ctx) error {

	user := new(UserModel)

	user.Username = c.Params("username")

	if user.Find() {
		log.Printf("Find user")
		return c.Status(401).JSON(fiber.Map{
			"message": "Find user",
		})
	}

	user.Roles = 1 << 0
	user.Create()

	options, sessionData, err := web.BeginRegistration(*user)

	if err != nil {
		log.Print(err)
		return c.SendStatus(401)
	}

	session := new(UserSessions)
	session.displayName = options.Response.User.Name
	session.sessionData = sessionData
	session.expiration = 60 * 60 * 3600

	go session.deleteAfter()
	sessions[options.Response.User.Name] = session

	return c.JSON(fiber.Map{
		"Options": options,
	})

}

func RegisterEnd(c *fiber.Ctx) error {
	user := new(UserModel)
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

	session, ok := sessions[user.Username]

	if ok == false {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"message": "session Not exist",
		})
	}
	creds, err := web.CreateCredential(user, *session.sessionData, credential)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}

	user.credentals = append(user.credentals, *creds)

	user.saveCredentials()

	session.sessionCred = creds
	session.expiration = 24 * 3600 * 2

	token, err := createJWT(session.displayName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	session.jwt = token
	go session.deleteAfter()
	sessions[user.Username] = session

	return c.JSON(fiber.Map{
		"token": session.jwt,
	})
}

func RegisterPassword(c *fiber.Ctx) error {
	user := new(UserModel)
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

	session := new(UserSessions)

	session.displayName = user.Username
	session.expiration = 24 * 3600 * 2

	token, err := createJWT(session.displayName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	session.jwt = token
	sessions[user.Username] = session
	go session.deleteAfter()

	user.Create()

	return c.JSON(fiber.Map{
		"token": session.jwt,
	})

}
