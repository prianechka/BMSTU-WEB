package jwtUtils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"src/objects"
)

var (
	PrivateKey = []byte("secret key")
)

func CreateJWTToken(login string, role objects.Levels) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"role":  role,
	})

	tokenString, _ := token.SignedString(PrivateKey)
	return tokenString
}

func GetRoleFromJWT(inToken string) (result objects.Levels) {
	result = objects.NonAuth
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return PrivateKey, nil
	}

	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err == nil && token.Valid {
		payload, ok := token.Claims.(jwt.MapClaims)
		if ok {
			floatRole := payload["role"].(float64)
			result = objects.Levels(floatRole)
		}
	}
	return result
}
