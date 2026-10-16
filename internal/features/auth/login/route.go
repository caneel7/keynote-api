package login

import (
	"context"
	"net/http"
	"time"

	httpError "github.com/caneel7/keynote-api/internal/errors"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
	"github.com/danielgtaylor/huma/v2"
)

type Authenticate interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (sessions.TokensAndSession, error)
}

type LoginRequest struct {
	Body struct {
		Email    string `json:"email" format:"email" maxLength:"254"`
		Password string `json:"password"`
	}
}

type LoginResponse struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
	Body struct {
		Token     string `json:"token"`
		TokenType string `json:"token_type" example:"Bearer"`
		ExpiresIn int    `json:"expires_in" doc:"Access token lifetime in seconds"`
	}
}

func LoginRoute(api huma.API, svc Authenticate) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "login",
			Method:        http.MethodPost,
			Path:          "/auth/login",
			DefaultStatus: http.StatusOK,
			Tags:          []string{"Auth"},
			Summary:       "User Login",
			Description:   "User login endpoint, return jwt",
		},
		func(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
			r, err := svc.Login(ctx,req.Body.Email,req.Body.Password)
			if err != nil {
				return nil, httpError.From(err)
			}

			res := &LoginResponse{}
			res.Body.Token = r.AccessToken
			res.Body.TokenType = "Bearer"
			res.Body.ExpiresIn = int(time.Until(r.ExpireAT).Seconds())
			res.SetCookie = http.Cookie{
				Name: "refresh_token",
				Value: r.RefreshToken,
				Path: "/auth",
				Expires: time.Now().Add(r.RefreshTokenExpireAt),
				HttpOnly: true,
				Secure: true,
				SameSite: http.SameSiteStrictMode,
			}

			return res,nil
		},
	)
}
