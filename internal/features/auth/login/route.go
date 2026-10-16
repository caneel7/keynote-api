package login

import (
	"context"
	"net/http"

	httpError "github.com/caneel7/keynote-api/internal/errors"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
	"github.com/danielgtaylor/huma/v2"
)

type Authenticate interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (sessions.Tokens, error)
}

type LoginRequest struct {
	Body struct {
		Email    string `json:"email" format:"email" maxLength:"254"`
		Password string `json:"password"`
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
		func(ctx context.Context, req *LoginRequest) (*sessions.Response, error) {
			t, err := svc.Login(ctx, req.Body.Email, req.Body.Password)
			if err != nil {
				return nil, httpError.From(err)
			}
			return sessions.NewResponse(t), nil
		},
	)
}
