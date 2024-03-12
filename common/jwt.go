package common

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(
	userId string,
	sessionId string,
	secret string,
	expired time.Time,
) string {
	token := jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		&jwt.StandardClaims{ExpiresAt: expired.Unix(), Id: sessionId, Subject: userId},
	)
	tokenString, _ := token.SignedString([]byte(secret))

	return tokenString
}

func ValidateToken(accessToken string, secret string) (*jwt.StandardClaims, bool) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, false
	}

	claims := token.Claims.(*jwt.StandardClaims)

	return claims, true
}
