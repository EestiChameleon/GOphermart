package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/EestiChameleon/GOphermart/internal/app/cfg"
	"github.com/robbert229/jwt"
)

var (
	ErrInvalidToken = errors.New("failed to decode the provided Token")
)

func EncryptPass(pass string) string {
	h := md5.New()
	h.Write([]byte(pass))
	return hex.EncodeToString(h.Sum(nil))
}

func JWTEncodeUserID(value interface{}) (string, error) {
	return JWTEncode("sub", value)
}

func JWTEncode(key string, value interface{}) (string, error) {
	algorithm := jwt.HmacSha256(cfg.Envs.CryptoKey)

	claims := jwt.NewClaim()
	claims.Set(key, value)

	token, err := algorithm.Encode(claims)
	if err != nil {
		return ``, err
	}

	if err = algorithm.Validate(token); err != nil {
		return ``, err
	}

	return token, nil
}

func JWTDecodeUserID(token string) (int, error) {
	value, err := JWTDecode(token, "sub")
	if err != nil {
		return -1, err
	}
	return int(value.(float64)), nil
}

func JWTDecode(token, key string) (interface{}, error) {
	algorithm := jwt.HmacSha256(cfg.Envs.CryptoKey)

	if err := algorithm.Validate(token); err != nil {
		return nil, err
	}

	claims, err := algorithm.Decode(token)
	if err != nil {
		return nil, err
	}

	return claims.Get(key)
}
