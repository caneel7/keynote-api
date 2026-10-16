package refresh

import (
	"context"
	"net/http"

	"github.com/caneel7/keynote-api/internal/errors"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
	"github.com/danielgtaylor/huma/v2"
)

type Refresh interface {
	RefreshToken(ctx context.Context, token string) (sessions.Tokens, error)
}

type RefreshTokenRequest struct {
	RefreshToken string `cookie:"refresh_token"`
}

func RefreshTokenRoute(api huma.API, svc Refresh) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "refreshtoken",
			Method:        http.MethodPost,
			Path:          "/auth/refresh",
			DefaultStatus: http.StatusOK,
			Tags:          []string{"Auth"},
			Summary:       "Refresh token",
			Description:   "Rotates the refresh token cookie and returns a new access token",
		},
		func(ctx context.Context, r *RefreshTokenRequest) (*sessions.Response, error) {
			t, err := svc.RefreshToken(ctx, r.RefreshToken)
			if err != nil {
				return nil, errors.From(err)
			}
			return sessions.NewResponse(t), nil
		},
	)
}
