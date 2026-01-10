package jwt

import (
	"errors"
	"time"

	"github.com/canonflow/canonflow-go-ddd/internal/domain/user/model"
	"github.com/gin-gonic/gin"
	jwtLibrary "github.com/golang-jwt/jwt/v5"
)

var (
	TOKEN_COOKIE   = "Authorization"
	TOKEN_KEY      = "TOKEN"
	USER_KEY       = "USER"
	PATH           = ""
	DOMAIN         = ""
	SECURE         = false
	HTTP_ONLY      = true
	TOKEN_DURATION = time.Hour * 6
)

/*
	JWT DATA
	{
		"sub": user_id,
		"username": username,
		"exp": expiration_time
	}
*/

func CreateToken(user *model.User, secretkey string) (string, error) {
	// TODO: Generate JWT Token
	token := jwtLibrary.NewWithClaims(jwtLibrary.SigningMethodHS256, jwtLibrary.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(TOKEN_DURATION).Unix(),
	})

	//* Sign and Get the complete encoded token as a string using the secret key
	tokenString, err := token.SignedString([]byte(secretkey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DeleteToken(ctx *gin.Context) {
	ctx.SetCookie(TOKEN_COOKIE, "delete", -1, DOMAIN, DOMAIN, SECURE, HTTP_ONLY)
}

func GetToken(ctx *gin.Context) (string, error) {
	token, err := ctx.Cookie(TOKEN_COOKIE)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(tokenString string, secret string) (jwtLibrary.MapClaims, error) {
	token, err := jwtLibrary.Parse(tokenString, func(token *jwtLibrary.Token) (any, error) {
		return []byte(secret), nil
	}, jwtLibrary.WithValidMethods([]string{jwtLibrary.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwtLibrary.MapClaims)

	if !ok {
		return nil, errors.New("Cannot claim JWT Token")
	}

	return claims, nil
}
