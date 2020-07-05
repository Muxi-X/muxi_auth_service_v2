package oauth

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ResolveAccessToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"]
		if !ok {
			return "", errors.New("Token not include `cap` field.")
		}
		return sub.(string), nil
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

func ParseRequest(c *gin.Context) (uint64, error) {
	token := c.GetHeader("token")
	if token == "" {
		return 0, errors.New("token is required")
	}
	return GetUserIDFromToken(token)
}
