package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
)

var jwtKey = os.Getenv("JWT_KEY")

func init() {
	if jwtKey == "" {
		jwtKey = "hello@Yrd!Kplcloud"
	}
}

func JwtKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	} else {
		return []byte(GetJwtKey()), nil
	}
}

func GetJwtKey() string {
	return jwtKey
}
