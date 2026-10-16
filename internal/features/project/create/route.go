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

type ProjectCreator interface {
	CreateProject(
		ctx context.Context,
		name string,
		tenantID uuid.UUID,
		workspaceID uuid.UUID,
		userID uuid.UUID,
	) (gen.Project, error)
}

type CreateProjectRequest struct {
	TenantID string `header:"X-Tenant-Id" required:"true" format:"uuid"`
	WorkspaceID string `header:"X-Workspace-Id" required:"true" format:"uuid"`
	Body struct {
		Name string `json:"name" minLength:"3" maxLength:"128"`
	}
}

type CreateProjectResponse struct {
	Body struct {
		ID string `json:"id" format:"uuid"`
		Name string `json:"name"`
		TenantID string `json:"tenant_id" format:"uuid"`
		WorkspaceID string `json:"workspace_id" format:"uuid"`
	}
}

func CreateProjectRoute(api huma.API,svc ProjectCreator) {
	huma.Register(
		api,
		huma.Operation{
			OperationID: "create-project",
			Method: http.MethodPost,
			Path: "/project",
			DefaultStatus: http.StatusCreated,
			Tags: []string{"Project"},
			Summary: "Create Project",
			Description: "Create new project under a workspace",
		},
		func (ctx context.Context,req *CreateProjectRequest) (*CreateProjectResponse,error) {

			userID, err := auth.GetUserIDFromClaims(ctx)
			if err != nil {
				return &CreateProjectResponse{},errors.From(err)
			}

			tenantID, workspaceID := req.TenantID, req.WorkspaceID
			tenantUUID, err := uuid.Parse(tenantID)
			if err != nil {
				return &CreateProjectResponse{}, errors.From(err)
			}
			workspaceUUID, err := uuid.Parse(workspaceID)
			if err != nil {
				return &CreateProjectResponse{}, errors.From(err)
			}

			p, err := svc.CreateProject(
				ctx,
				req.Body.Name,
				tenantUUID,
				workspaceUUID,
				userID,
			)
			if err != nil {
				return &CreateProjectResponse{}, errors.From(err)
			}

			response := &CreateProjectResponse{}
			response.Body.ID = p.ID.String()
			response.Body.Name = p.Name
			response.Body.TenantID = p.TenantID.String()
			response.Body.WorkspaceID = p.WorkspaceID.String()
			return response, nil
		},
	)
}