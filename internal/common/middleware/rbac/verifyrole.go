package rbac

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/caneel7/keynote-api/internal/constants"
	sqlc "github.com/caneel7/keynote-api/internal/db/sqlc"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/google/uuid"

	"github.com/danielgtaylor/huma/v2"
)

type Authorizer struct {
	ctx  context.Context
	db   *gen.Client
	sqlc *sqlc.Queries
	api  huma.API
}

func New(
	ctx context.Context,
	db *gen.Client,
	sqlc *sqlc.Queries,
	api huma.API,
) *Authorizer {
	return &Authorizer{ctx, db, sqlc, api}
}

const TenantContextKey = "TenantID"
const WorkspaceContextKey = "WorkspaceID"
const ProjectContextKey = "ProjectID"

func (a *Authorizer) Require(key string) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {

		keyModule := strings.Split(key, ":")[0]

		err := a.checkModule(keyModule, ctx)

		if err != nil {
			huma.WriteErr(a.api, ctx, http.StatusUnauthorized, "Invalid authorization header")
			return
		}



		next(ctx)
	}
}

func (a *Authorizer) checkModule(key string, ctx huma.Context) error {

	module, _ := constants.ResourceToScope[key]

	var headers = map[string]uuid.NullUUID{}

	for _, m := range module {
		rawHeader, parseErr := uuid.Parse(ctx.Header(m))
		if parseErr != nil || strings.TrimSpace(rawHeader.String()) == "" {
			return fmt.Errorf("cannot parse required header")
		}
		headers[m] = uuid.NullUUID{UUID: rawHeader, Valid: true}
	}

	validTenantId := headers["X-Tenant-Id"].Valid

	if !validTenantId {
		return fmt.Errorf("No tenant header")
	}

	ok, err := a.sqlc.VerifyProjectHierarchy(ctx.Context(), sqlc.VerifyProjectHierarchyParams{
		TenantID:    headers["X-Tenant-Id"].UUID,
		WorkspaceID: headers["X-Workspace-Id"],
		ProjectID:   headers["X-Project-Id"],
	})

	if err != nil || !ok {
		return fmt.Errorf("verify hierarchy: %w", err)
	}

	ctx = huma.WithValue(ctx, TenantContextKey, headers["X-Tenant-Id"].UUID)
	ctx = huma.WithValue(ctx, WorkspaceContextKey, headers["X-Workspace-Id"].UUID)
	ctx = huma.WithValue(ctx, ProjectContextKey, headers["X-Project-Id"].UUID)

	return nil
}

func (a *Authorizer) verifyPermission(key string,ctx huma.Context) error {


	

	return nil
}