package create

import (
	"context"
	"net/http"

	"github.com/caneel7/keynote-api/internal/common/middleware/auth"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/errors"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type TenantCreator interface {
	Create(
		ctx context.Context,
		name string,
		userID uuid.UUID,
	) (gen.Tenant, gen.Workspace, error)
}

type CreateTenantRequest struct {
	Body struct {
		Name string `json:"name" minLength:"3" maxLength:"128"`
	}
}

type CreateTenantResponse struct {
	Body struct {
		Tenant struct {
			ID   string `json:"id" format:"uuid"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"tenant"`
		Workspace struct {
			ID   string `json:"id" format:"uuid"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		}
	}
}

func CreateTenantRoute(api huma.API, svc TenantCreator) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "ceate-tenant",
			Method:      http.MethodPost,
			Path:        "/tenant",
			Tags:        []string{"Tenant"},
			Summary:     "Create Tenant",
			Description: "Create new tenant to be called when new user registers and chooses to create workplace",
		},
		func(ctx context.Context, r *CreateTenantRequest) (*CreateTenantResponse, error) {

			claims, _ := auth.ClaimFromContext(ctx)
			res := &CreateTenantResponse{}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				return &CreateTenantResponse{}, err
			}

			t, w, err := svc.Create(
				ctx,
				r.Body.Name,
				userID,
			)

			if err != nil {
				return res, errors.From(err)
			}

			res.Body.Tenant.ID = t.ID.String()
			res.Body.Tenant.Name = t.Name
			res.Body.Tenant.Slug = t.Slug

			res.Body.Workspace.ID = w.ID.String()
			res.Body.Workspace.Name = w.Name
			res.Body.Workspace.Slug = w.Slug

			return res, nil
		},
	)
}
