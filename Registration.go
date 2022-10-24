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

	if user.Find(db) {
		log.Printf("Find user")
		return c.Status(401).JSON(fiber.Map{
			"message": "Find user",
		})
	}

	user.Roles = 1 << 0
	user.Create(db)

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

	if !user.Find(db) {
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

	user.saveCredentials(db)

	session.sessionCred = creds
	session.expiration = 24 * 3600 * 2
	go session.deleteAfter()

	return c.JSON(fiber.Map{
		"creds": creds,
	})
}
