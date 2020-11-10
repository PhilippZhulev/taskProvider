package handler

import (
	"encoding/base64"
	"log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"gitlab.com/taskProvider/services/user/internal/app/configurator"
	"gitlab.com/taskProvider/services/user/internal/app/store/sqlstore"
)

// HandleLogin ...
// Проверить логин и пароль
func (init Init) HandleLogin(in string, store *sqlstore.Store, config *configurator.Config) ([]string, error) {
	decoded, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}

	result := strings.Split(string(decoded), ":")
	u, err := store.User().FindByKey("login_name", result[0])
	if err != nil {
		log.Println(err)
		return nil, errIncorrectEmailOrPassword;
	}

	if !init.hesh.CheckPasswordHash(result[1], u.EncryptedPassword, config.Salt) {
		return nil, errIncorrectEmailOrPassword;
	}

	return []string{u.Email, u.UUID}, nil
}

// HandleCreateToken ...
// Стенерировать токен
func (init Init) HandleCreateToken(attr []string, tokenAuth *jwtauth.JWTAuth) (string, error) {
	claims := jwt.MapClaims{"user_id": attr[0], "uuid": attr[1]}
	jwtauth.SetExpiryIn(claims, 1000*time.Second)
	_, tokenString, err := tokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
