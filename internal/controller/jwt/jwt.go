package controller

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"

	ierr "webChat/internal/errors"
)

func GenerateToken(secretKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseToken(accessToken string, signingKey []byte) error {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ierr.Internal(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return signingKey, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return ierr.ErrInvalidAccessToken
	}

	return ierr.ErrInvalidAccessToken
}
