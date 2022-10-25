package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserSessions struct {
	sessionData *webauthn.SessionData
	sessionCred *webauthn.Credential
	displayName string
	jwt         string
	expiration  uint64
}

var (
	sessions map[string]*UserSessions
	web      *webauthn.WebAuthn
	db       *gorm.DB
	err      error
)

/*
*  return true only if a session contains AAGUID
 */
func checkAuthn(c *fiber.Ctx) *UserSessions {
	value, ok := c.GetReqHeaders()["Authorization"]
	if ok == false {
		return nil
	}
	authType := strings.Split(value, " ")
	if authType[0] != "Bearer" || len(authType) < 2 {
		return nil
	}

	auth := strings.Split(authType[1], "?")

	if len(auth) < 2 {
		return nil
	}

	for _, v := range sessions {

		if strings.Compare(auth[1], base64.URLEncoding.EncodeToString(v.sessionCred.Authenticator.AAGUID)) == 0 {
			for _, v2 := range v.sessionData.AllowedCredentialIDs {
				if strings.Compare(base64.URLEncoding.EncodeToString(v2), strings.Replace(auth[0], "/", "_", 1)) == 0 {
					return v
				}

			}
			return nil
		}

	}
	return nil
}

func main() {
	/* env vars */
	godotenv.Load(".env")
	postgresHost := os.Getenv("PostgresHost")
	postgresUser := os.Getenv("PostgresUser")
	postgresPassword := os.Getenv("PostgresPassword")
	postgresDatabase := os.Getenv("PostgresDatabase")
	postgresPort := os.Getenv("postgresPort")
	RPDiplayName := os.Getenv("RPDisplayName")
	RPID := os.Getenv("RPID")
	ROrigin := os.Getenv("RPOrigin")
	RPIcon := os.Getenv("RPIcon")
	appListen := os.Getenv("AppListen")

	app := fiber.New()
	sessions = make(map[string]*UserSessions)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// db Initialisaiton
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", postgresHost, postgresUser, postgresPassword, postgresDatabase, postgresPort)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&UserModel{}, &RoleModel{})

	// webauthn init

	web, err = webauthn.New(&webauthn.Config{
		RPDisplayName: RPDiplayName, // Display Name for your site
		RPID:          RPID,         // Generally the FQDN for your site
		RPOrigin:      ROrigin,      // The origin URL for WebAuthn requests
		RPIcon:        RPIcon,       // Optional icon URL for your site
	})
	if err != nil {
		fmt.Println(err)
	}

	//app routes
	app.Post("register/start/:username", RegistrationStart)

	app.Post("register/end/:username", RegisterEnd)

	app.Post("login/start/:username", LoginStart)

	app.Post("login/end/:username", LoginEnd)

	api := app.Group("/api", func(c *fiber.Ctx) error {
		if checkAuthn(c) == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()

	})

	user := api.Group("user")
	user.Get("/", func(c *fiber.Ctx) error {
		user := new(UserModel)
		userSession := checkAuthn(c)
		user.Username = userSession.displayName
		return c.Status(200).JSON(user.Get())
	})

	user.Get("/logout", func(c *fiber.Ctx) error {
		userSession := checkAuthn(c)
		delete(sessions, userSession.displayName)
		return c.Status(200).JSON(fiber.Map{
			"message": "logout",
		})
	})

	user.Patch("/", func(c *fiber.Ctx) error {
		user := new(UserModel)
		if err := c.BodyParser(user); err != nil {
			fmt.Println("error = ", err)
			return c.SendStatus(200)
		}
		userSession := checkAuthn(c)
		user.Username = userSession.displayName

		user.Update()

		return c.Status(200).JSON(user)

	})

	user.Delete("/", func(c *fiber.Ctx) error {
		user := new(UserModel)
		userSession := checkAuthn(c)
		user.Username = userSession.displayName

		user.Delete()
		delete(sessions, user.Username)

		return c.JSON(fiber.Map{
			"message": "deleted",
		})
	})
	user.Delete("/cred", func(c *fiber.Ctx) error {
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

	//app run
	log.Fatal(app.Listen(appListen))

}
