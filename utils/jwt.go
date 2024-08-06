package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("KGSfcuc$hehfejn#fjebf&cJYRTYfj")

type Claims struct {
	EmailId string `json:"emailId"`
	jwt.RegisteredClaims
}

func GenerateJWT(emailId string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		EmailId: emailId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
