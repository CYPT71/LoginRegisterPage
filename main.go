package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	// "github.com/joho/godotenv"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	sessions map[string]*UserSessions
	web      *webauthn.WebAuthn
	db       *gorm.DB
	err      error
)

func main() {
	/* env vars */
	if _, err := os.Stat(".env"); err == nil {
		godotenv.Load(".env")
	}

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

	db.AutoMigrate(&UserModel{})

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

	app.Get("/checkUser/:username", CheckUserName)
	//app routes
	app.Post("register/start/:username", RegistrationStart)

	app.Post("register/end/:username", RegisterEnd)

	app.Post("register/password/:username", RegisterPassword)

	app.Post("login/start/:username", LoginStart)

	app.Post("login/end/:username", LoginEnd)

	app.Post("login/password/:username", loginPassword)

	UserBootstrap(app.Group("user", func(c *fiber.Ctx) error {

		if checkAuthn(c) == nil {
			log.Println(c.GetReqHeaders())
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	}))

	//app run
	log.Fatal(app.Listen(appListen))

}
