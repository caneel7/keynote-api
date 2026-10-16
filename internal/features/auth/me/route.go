package me

import (
	"context"
	"net/http"

	"github.com/caneel7/keynote-api/internal/common/middleware/auth"
	"github.com/caneel7/keynote-api/internal/errors"
	"github.com/caneel7/keynote-api/internal/types"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type GetMe interface {
	GetMe(ctx context.Context, userID uuid.UUID) (types.User, error)
}

type GetMeResponse struct {
	Body struct {
		User struct {
			ID    string `json:"id" format:"uuid"`
			Email string `json:"email"`
		} `json:"user" doc:"Authenticated user"`
	}
}

func GetMeRoute(api huma.API, svc GetMe) {
	huma.Register(
		api,
		huma.Operation{
			OperationID:   "getme",
			Method:        http.MethodGet,
			Path:          "/auth/me",
			DefaultStatus: http.StatusOK,
			Tags:          []string{"Auth"},
			Summary:       "Me (Authenticated user)",
			Description:   "Returns authenticated user via jwt",
		},
		func(ctx context.Context, input *struct{}) (*GetMeResponse, error) {
			userID, err := auth.GetUserIDFromClaims(ctx)
			if err != nil {
				return &GetMeResponse{}, errors.From(err)
			}

			user, err := svc.GetMe(ctx, userID)
			if err != nil {
				return &GetMeResponse{}, errors.From(err)
			}

			res := &GetMeResponse{}
			res.Body.User.ID = user.ID.String()
			res.Body.User.Email = user.Email
			return res, nil
		},
	)
}
