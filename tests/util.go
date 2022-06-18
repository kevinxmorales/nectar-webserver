package tests

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
)

const BaseUrl = "http://localhost:8080/api"
const Version = "v1"
const plantEndpoint = "plant"
const userEndpoint = "user"

func GetToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}
