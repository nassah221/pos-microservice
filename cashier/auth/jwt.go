package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Authenticator interface {
	IssueToken(user, id string, duration time.Duration) (string, error)
	VerifyToken(token string) (*jwt.Token, error)
}

type auth struct {
	key string
}

func NewAuthenticator(key string) Authenticator {
	return &auth{key}
}

func (a *auth) IssueToken(user, id string, duration time.Duration) (string, error) {
	now := time.Now()

	defaultExpiry := time.Second * 1
	if duration == 0 {
		defaultExpiry = time.Hour * 24
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// standard claims
		"aud": "api",
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"exp": now.Add(defaultExpiry).Unix(),
		"iss": "cashier",

		// custom claims
		"user": user,
		"id":   id,
	})

	tokenString, err := token.SignedString([]byte(a.key))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func (a *auth) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(a.key), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	return token, nil
}
