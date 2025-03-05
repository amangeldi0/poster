package auth

import (
	"github.com/google/uuid"
	assert2 "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteCookie(t *testing.T) {
	assert := assert2.New(t)
	cookieKey := "test_key"

	writer := httptest.NewRecorder()
	DeleteCookie(cookieKey, writer)

	resp := writer.Result()
	cookies := resp.Cookies()

	assert.Len(cookies, 1, "Must have one cookie")
	cookie := cookies[0]

	assert.Equal(cookieKey, cookie.Name, "cookie name should be equal")
	assert.Equal("", cookie.Value, "Cookie value should be empty")
	assert.Equal("/", cookie.Path, "Cookie have to set on '/'")
	assert.True(cookie.HttpOnly, "Cookie have to be HttpOnly")
	assert.True(cookie.Secure, "Cookie have to be Secure")
	assert.Equal(http.SameSiteStrictMode, cookie.SameSite, "SameSite have to be Strict")
	assert.Equal(-1, cookie.MaxAge, "MaxAge have to be cookie")
}

func TestHashPassword(t *testing.T) {
	assert := assert2.New(t)
	password := "password"

	res, err := HashPassword(password)

	assert.Nil(err)
	assert.NotEqual(password, res)
}

func TestCheckPasswordHash(t *testing.T) {
	assert := assert2.New(t)
	password := "password"

	res, err := HashPassword(password)

	assert.Nil(err)
	assert.NotEqual(password, res)

	err = CheckPasswordHash(password, res)
	assert.Nil(err)
}

func TestGenerateCode(t *testing.T) {
	assert := assert2.New(t)
	code, err := GenerateCode()
	assert.Nil(err)
	assert.NotEqual("", code)
	assert.Len(code, 6)
}

func TestGenerateAccessToken(t *testing.T) {
	userUUID := uuid.New()
	assert := assert2.New(t)

	token, err := GenerateAccessToken(userUUID.String())

	assert.Nil(err)
	assert.NotEqual("", token)
}

func TestGenerateRefreshToken(t *testing.T) {
	userUUID := uuid.New()
	assert := assert2.New(t)

	token, err := GenerateAccessToken(userUUID.String())

	assert.Nil(err)
	assert.NotEqual("", token)
}

func TestVerifyToken(t *testing.T) {
	userUUID := uuid.New()
	assert := assert2.New(t)

	token, err := GenerateAccessToken(userUUID.String())

	assert.Nil(err)
	assert.NotEqual("", token)

	d, err := VerifyToken(token)

	assert.Nil(err)

	assert.Equal(userUUID.String(), d.UserID)
}
