package main

import (
	"log"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Username string `json:"username"`
	Ndf      int    `json:"ndf"`
	jwt.Claims
}

var sampleSecretKey = []byte(generateKey())

func generateKey() string {
	return "AZERTY"
}

func createJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"username": username,
		"ndf":      24 * 3600 * 2,
	})
	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func checkJWT(session *UserSessions, tokenString string) bool {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte("AllYourBase"), nil
		})

	if err != nil {
		log.Println(err.Error())

		return false
	}
	_, ok := token.Claims.(*Claims)
	return ok && token.Valid

}
