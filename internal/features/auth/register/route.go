package register

import (
	"context"
	"net/http"
	"time"

	httpError "github.com/caneel7/keynote-api/internal/errors"
	"github.com/caneel7/keynote-api/internal/features/auth/login"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
	"github.com/danielgtaylor/huma/v2"
)

type Registrar interface {
	Register(
		ctx context.Context,
		email string,
		password string,
	) (sessions.TokensAndSession, error)
}

type RegisterUserRequest struct {
	Body struct {
		Email    string `json:"email" format:"email" maxLength:"254"`
		Password string `json:"password" minLength:"8" maxLength:"128"`
	}
}


func RegisterUserRoute(api huma.API,svc Registrar) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "register",
			Method: http.MethodPost,
			Path: "/auth/register",
			DefaultStatus: http.StatusCreated,
			Tags: []string{"Auth"},
			Summary: "Register new user",
			Description: "Register new user in system",
		},
		func (ctx context.Context,req *RegisterUserRequest) (*login.LoginResponse,error) {
			r, err := svc.Register(ctx,req.Body.Email,req.Body.Password)
			if err != nil {
				return nil, httpError.From(err)
			}
			res := &login.LoginResponse{}
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