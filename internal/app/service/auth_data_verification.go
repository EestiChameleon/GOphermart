package service

import (
	"errors"
	"github.com/EestiChameleon/GOphermart/internal/app/service/methods"
	"github.com/EestiChameleon/GOphermart/internal/models"
)

var (
	ErrWrongAuthData = errors.New("wrong authentication data")
)

func CheckAuthData(ld models.LoginData) (string, error) {
	u := methods.NewUser(ld.Login, ld.Password)
	if err := u.GetByLogin(); err != nil {
		return "", err
	}

	if EncryptPass(ld.Password) != u.Password {
		return "", ErrWrongAuthData
	}

	return JWTEncodeUserID(u.ID)
}
