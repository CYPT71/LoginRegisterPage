package controllers

import (
	"strings"

	"webauthn_api/internal/domain"
	"webauthn_api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func UserBootstrap(app fiber.Router) {

	app.Get("/", about)

	app.Get("/all", allUsers)

	app.Get("/logout", logout)

	app.Patch("/", editUser)

	app.Delete("/", deleteUser)

	app.Delete("/cred", deleteCred)

}

// Get User
// @Summary Get about me
// @Description get all information about me
// @Tags Users
// @Success 200 {UserModel} domain.UserModel
// @Failure 404
// @Router /user [get]
func about(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	userSession := utils.CheckAuthn(c)
	if userSession == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user.Username = userSession.DisplayName
	return c.Status(200).JSON(user.Get())
}

// Logout
// @Summary Just Logout
// @Tags Users
// @Success 200 {array} domain.UserModel
// @Unthorized 401
// @Failure 500 nil object
// @Router /user/all
func allUsers(c *fiber.Ctx) error {

	userSession := utils.CheckAuthn(c)
	if userSession == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user := new(domain.UserModel)
	user.Username = userSession.DisplayName

	if user == nil || user.Get().Permission&domain.Permissions["owner"] != 1 {
		return c.SendStatus(401)
	}

	users, err := domain.GetAllUsers()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.Status(200).JSON(users)
}

// Logout
// @Summary Just Logout
// @Tags Users
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user/logout [get]
func logout(c *fiber.Ctx) error {
	userSession := utils.CheckAuthn(c)
	if userSession == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	delete(utils.Sessions, userSession.DisplayName)
	return c.Status(200).JSON(fiber.Map{
		"message": "logout",
	})
}

// Edit me
// @Summary  edit user
// @Tags Users
// @Description edit user information
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user [patch]
func editUser(c *fiber.Ctx) error {

	userSession := utils.CheckAuthn(c)
	if userSession == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userIn := new(domain.UserModel)
	userIn.Username = userSession.DisplayName
	userIn = userIn.Get()

	user := new(utils.PartialUser)
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

}

// Delete me
// @Summary  delete account
// @Tags Users
// @Description delete user account
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user [delete]
func deleteUser(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	userSession := utils.CheckAuthn(c)
	if userSession == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user.Username = userSession.DisplayName

	user.Delete()
	delete(utils.Sessions, user.Username)

	return c.JSON(fiber.Map{
		"message": "deleted",
	})
}

// Delete credential
// @Summary  delete credential
// @Tags Users
// @Description delete webauthn credential
// @Success 200 {array} domain.UserModel
// @Failure 404 nil object
// @Router /user/cred [delete]
func deleteCred(c *fiber.Ctx) error {
	user := new(domain.UserModel)
	userSession := utils.CheckAuthn(c)
	if userSession == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user.Username = userSession.DisplayName
	user = user.Get()
	if user == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user.Incredentials = strings.Split(user.Incredentials, ";")[0]
	user.Update()

	return c.Status(200).JSON(user)
}
