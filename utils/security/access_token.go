package security

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/config"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model"
	"log"
	"time"
)

type AccessToken interface {
	CreateAccessToken(cred *model.UserCredential) (string, error)
	VerifyAccessToken(tokenString string) (jwt.MapClaims, error)
}

type accessToken struct {
	cfg config.TokenConfig
}

func NewAccessToken(config config.TokenConfig) AccessToken {
	return &accessToken{
		cfg: config,
	}
}

func (t *accessToken) CreateAccessToken(cred *model.UserCredential) (string, error) {
	now := time.Now().UTC()
	end := now.Add(t.cfg.AccessTokenLifeTime)
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer: t.cfg.ApplicationName,
		},
		Username: cred.UserName,
		Email:    cred.UserName,
	}
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = end.Unix()
	token := jwt.NewWithClaims(
		t.cfg.JwtSigningMethod,
		claims,
	)
	return token.SignedString([]byte(t.cfg.JwtSignatureKey))
}

func (t *accessToken) VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != t.cfg.JwtSigningMethod {
			return nil, fmt.Errorf("signing method invalid")
		}

		return []byte(t.cfg.JwtSignatureKey), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims["iss"] != t.cfg.ApplicationName {
		log.Println("Token Invalid")
		return nil, err
	}
	return claims, nil
}
