package main

import (
	"log"

	"math/rand"

	"github.com/golang-jwt/jwt"
)

var sampleSecretKey = []byte(generateKey(20))

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateKey(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func createJWT(session UserSessions) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  session.displayName,
		"AuthToken": string(session.sessionData.UserID),
	})
	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func checkJWT(session *UserSessions, tokenString string) bool {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	if err != nil {
		log.Println(err.Error())

		return false
	}
	if token.Method != jwt.SigningMethodHS256 {
		log.Println(token.Valid)
		return false
	}

	i := 1 << 0
	for _, val := range claims {
		if val == session.displayName {
			i |= 1 << 1
		}
		if val == string(session.sessionData.UserID) {
			i |= 1 << 0
		}
	}

	return i == 3

}
