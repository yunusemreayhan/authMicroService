package key

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yunusemreayhan/goAuthMicroService/internal/config"
	"github.com/yunusemreayhan/goAuthMicroService/internal/model"
	"time"
)

func CreateClaims(identity model.JWTIdentity) jwt.Claims {
	claims := jwt.MapClaims{
		"identifier": model.JWTIdentity{Id: identity.Id, Username: identity.Username, Email: identity.Email},
		"admin":      true,
		"exp":        time.Now().Add(time.Duration(config.GetConfig().TokenTimeoutSeconds) * time.Second).Unix(),
	}
	return claims
}

func CreateJwtToken(claims jwt.Claims, privateKey *rsa.PrivateKey) (string, error) {
	// todo either consider signing with user password or keep jwt duration in config short
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func CreateJwtTokenWithIdentity(identity model.JWTIdentity, privateKey *rsa.PrivateKey) (string, error) {
	claims := CreateClaims(identity)
	return CreateJwtToken(claims, privateKey)
}
