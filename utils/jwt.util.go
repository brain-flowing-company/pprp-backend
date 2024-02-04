package utils

import (
	"time"

	"github.com/brain-flowing-company/pprp-backend/internal/models"
	"github.com/golang-jwt/jwt"
)

func CreateJwtToken(session models.Session, maxAge time.Duration, jwtSecret string) (string, error) {
	issuedTime := time.Now().Unix()
	expiresTime := time.Now().Add(maxAge).Unix()

	customClaim := models.SessionClaim{
		Session: session,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedTime,
			ExpiresAt: expiresTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaim)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
