package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"

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
)

func main() {

	app := fiber.New()
	sessions = make(map[string]*UserSessions)

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	var err error

	// db Initialisaiton
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&UserModel{}, &RoleModel{})

	// webauthn init

	web, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Duo Labs",              // Display Name for your site
		RPID:          "localhost",             // Generally the FQDN for your site
		RPOrigin:      "http://localhost:8080", // The origin URL for WebAuthn requests

		RPIcon: "https://duo.com/logo.png", // Optional icon URL for your site
	})
	if err != nil {
		fmt.Println(err)
	}

	//app routes
	app.Post("register/start/:username", RegistrationStart)

	app.Post("register/end/:username", RegisterEnd)

	app.Post("login/start/:username", LoginStart)

	app.Post("login/end/:username", LoginEnd)

	//app run
	log.Fatal(app.Listen(":3000"))

}
