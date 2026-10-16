package refresh

import (
	"context"
	"net/http"
	"time"

	"github.com/caneel7/keynote-api/internal/errors"
	login "github.com/caneel7/keynote-api/internal/features/auth/login"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
	"github.com/danielgtaylor/huma/v2"
)

type Refresh interface {
	RefreshToken(ctx context.Context,token string) (sessions.TokensAndSession,error)
}

type RefreshTokenRequest struct {
	RefreshToken string `cookie:"refresh_token"`
}

func RefreshTokenRoute(api huma.API,svc Refresh) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "refreshtoken",
			Method: http.MethodGet,
			Path: "/auth/refresh",
			DefaultStatus: http.StatusOK,
			Tags: []string{"Auth"},
			Summary: "Refresh token",
			Description: "Returns new access token",
		},
		func (ctx context.Context,r *RefreshTokenRequest) (*login.LoginResponse,error) {
			res := &login.LoginResponse{}

			out, err := svc.RefreshToken(ctx,r.RefreshToken)
			if err != nil {
				return res,errors.From(err)
			}

			res.SetCookie = http.Cookie{
				Name:     "refresh_token",
				Value:    out.RefreshToken,
				Path:     "/auth",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			res.Body.Token = out.AccessToken
			res.Body.TokenType = "Bearer"
			res.Body.ExpiresIn = int(time.Until(out.ExpireAT).Seconds())

			return res,nil
		},
	)
}