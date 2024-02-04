package models

import "github.com/golang-jwt/jwt"

type SessionClaim struct {
	jwt.StandardClaims
	Session Session `json:"session"`
}
