package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"webauthn_api/internal/domain"
	"webauthn_api/internal/http"
	"webauthn_api/internal/utils"

	"github.com/joho/godotenv"

	// "github.com/joho/godotenv"

	_ "webauthn_api/docs"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var err error

func main() {
	/* env vars */
	if _, err := os.Stat(".env"); err == nil {
		godotenv.Load(".env")
	}

	postgresHost := os.Getenv("PostgresHost")
	postgresUser := os.Getenv("PostgresUser")
	postgresPassword := os.Getenv("PostgresPassword")
	postgresDatabase := os.Getenv("PostgresDatabase")
	postgresPort := os.Getenv("PostgresPort")
	RPDiplayName := os.Getenv("RPDisplayName")
	RPID := os.Getenv("RPID")
	ROrigin := strings.Split(os.Getenv("RPOrigin"), ", ")
	// RPIcon := os.Getenv("RPIcon")
	appListen := os.Getenv("AppListen")

	utils.Sessions = make(map[string]*domain.UserSessions)

	// db Initialisaiton
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", postgresHost, postgresUser, postgresPassword, postgresDatabase, postgresPort)
	domain.Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	domain.Db.AutoMigrate(&domain.UserModel{})

	// webauthn init

	utils.Web, err = webauthn.New(&webauthn.Config{
		RPDisplayName: RPDiplayName, // Display Name for your site
		RPID:          RPID,         // Generally the FQDN for your site
		RPOrigins:     ROrigin,      // The origin URL for WebAuthn requests
		// :        RPIcon,       // Optional icon URL for your site
	})
	if err != nil {
		fmt.Println(err)
	}

	//app run
	log.Fatal(http.Http().ListenTLS(appListen, "./certs/server.crt", "./certs/server.key"))

}
