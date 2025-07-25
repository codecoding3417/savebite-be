package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"savebite/internal/domain/env"
	"savebite/pkg/log"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	ID    uuid.UUID
	Name  string
	Email string
}

type CustomJWTItf interface {
	Create(id uuid.UUID, name, email string) (string, error)
	Decode(tokenString string, claim *Claims) error
}

type CustomJWTStruct struct {
	secretKey   string
	expiredTime time.Duration
}

var JWT = getJWT()

func getJWT() CustomJWTItf {
	return &CustomJWTStruct{
		secretKey:   env.AppEnv.JwtSecretKey,
		expiredTime: env.AppEnv.JwtExpTime,
	}
}

func (j *CustomJWTStruct) Create(id uuid.UUID, name, email string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiredTime)),
		},
		Name:  name,
		Email: email,
		ID:    id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[JWT][Create] Failed to generate jwt token")
		return "", err
	}

	return tokenString, nil
}

func (j *CustomJWTStruct) Decode(tokenString string, claim *Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claim, func(_ *jwt.Token) (any, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}
