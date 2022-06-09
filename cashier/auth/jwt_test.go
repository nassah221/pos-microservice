package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

var jwtMaker Authenticator

func init() {
	jwtMaker = NewAuthenticator("secret")
}

func TestAuth(t *testing.T) {
	tokenString, err := jwtMaker.IssueToken("user", "id", time.Second)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwtMaker.VerifyToken(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, token)

	user := token.Claims.(jwt.MapClaims)["user"]
	id := token.Claims.(jwt.MapClaims)["id"]
	assert.Equal(t, "user", user)
	assert.Equal(t, "id", id)

	time.Sleep(time.Second * 3)
	token, err = jwtMaker.VerifyToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, token)
	assert.Equal(t, err.Error(), "failed to parse token: Token is expired")
}
