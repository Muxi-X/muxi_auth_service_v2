package captcha

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
	"time"
)

func GetCaptcha(length int) string {
	rand.Seed(time.Now().UnixNano())
	tempString := "%0" + strconv.Itoa(length) + "d"
	return fmt.Sprintf(tempString, rand.Int()%1000000)
}

func GenerateCaptchaToken(captchaCode string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["cap"] = captchaCode
	rtClaims["exp"] = time.Now().Add(time.Minute * 3).Unix()
	rt, err := refreshToken.SignedString([]byte(viper.GetString("SECRET")))
	if err != nil {
		return "", err
	}
	return rt, nil
}

func ResolveCaptchaToken(captchaToken string) (string, error) {
	token, err := jwt.Parse(captchaToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(viper.GetString("SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if cap, ok := claims["cap"]; !ok {
			return "", errors.New("Token not include `cap` field.")
		} else {
			return cap.(string), nil
		}
	}
	return "", errors.New("Unknown error.")
}
