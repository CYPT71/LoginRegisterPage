package main

import (
	"bytes"
	"encoding/base64"
	"log"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/gofiber/fiber/v2"
)

func LoginStart(c *fiber.Ctx) error {

	user := new(UserModel)
	user.Username = c.Params("username")
	if user.Find() == false {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "no user with this username",
		})
	}
	user.parseCredentials()

	options, sessionData, err := web.BeginLogin(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	session := new(UserSessions)

	session.sessionData = sessionData
	session.displayName = user.Username
	session.expiration = 60 * 3
	go session.deleteAfter()
	sessions[user.Username] = session

	return c.JSON(options)

}

func LoginEnd(c *fiber.Ctx) error {
	user := new(UserModel)
	user.Username = c.Params("username")
	if user.Find() == false {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "no user with this username",
		})
	}
	user.parseCredentials()

	session, ok := sessions[c.Params("username")]

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

	creds, err := web.ValidateLogin(user, *session.sessionData, parsedCredential)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": err,
		})
	}
	session.sessionCred = creds
	session.expiration = 24 * 3600 * 2
	go session.deleteAfter()

	sessions[user.Username] = session
	user.credentals = append(user.credentals, *creds)
	user.saveCredentials()

	return c.JSON(fiber.Map{
		"token": base64.URLEncoding.EncodeToString(creds.ID) + "?" + base64.URLEncoding.EncodeToString(creds.Authenticator.AAGUID),
	})
}
