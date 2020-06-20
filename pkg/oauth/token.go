package oauth

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

func ResolveAccessToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("oauth"), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"]; !ok {
			return "", errors.New("Token not include `cap` field.")
		} else {
			return sub.(string), nil
		}
	}
	return "", errors.New("Unknown error.")
}

func GetUserIDFromToken(token string) (uint64, error) {
	userID, err := ResolveAccessToken(token)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(userID, 10, 64)
}
