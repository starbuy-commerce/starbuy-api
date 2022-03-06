package jwt

import (
	"authentication-service/util"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(username string) string {

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["expiration"] = time.Now().Add(time.Hour * 24).Unix()
	claims["username"] = username
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	var config = util.GrabConfig()

	str, err := token.SignedString([]byte(config.Secret))

	if err != nil {
		log.Fatal(err)
	}

	return str
}