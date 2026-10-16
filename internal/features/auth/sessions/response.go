package sessions

import (
	"net/http"
	"time"
)

const refreshCookieName = "refresh_token"

// Response for common routes (register, login, refresh).
type Response struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body      struct {
		Token     string `json:"token"`
		TokenType string `json:"token_type" example:"Bearer"`
		ExpiresIn int    `json:"expires_in" doc:"Access token lifetime in seconds"`
	}
}

func NewResponse(t Tokens) *Response {
	res := &Response{}
	res.Body.Token = t.AccessToken
	res.Body.TokenType = "Bearer"
	res.Body.ExpiresIn = int(time.Until(t.AccessExpiresAt).Seconds())
	res.SetCookie = http.Cookie{
		Name:     refreshCookieName,
		Value:    t.RefreshToken,
		Path:     "/auth",
		Expires:  t.RefreshExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	return res
}
