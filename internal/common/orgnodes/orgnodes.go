package orgnodes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sqlc "github.com/caneel7/keynote-api/internal/db/sqlc"
	appError "github.com/caneel7/keynote-api/internal/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	TypeTenant    = "tenant"
	TypeWorkspace = "workspace"
	TypeProject   = "project"
)

type NewNode struct {
	ParentID  *uuid.UUID
	Type      string
	Name      string
	StatusID  uuid.UUID
	CreatedBy uuid.UUID
}

func AppendNode(ctx context.Context, q *sqlc.Queries, n NewNode) (sqlc.OrgNode, error) {
	id := uuid.New()

	var node sqlc.OrgNode
	var err error

	if n.ParentID == nil {
		if n.Type != TypeTenant {
			return sqlc.OrgNode{}, appError.New(appError.KindBadRequest, "Only a tenant can be created without a parent")
		}
		node, err = q.CreateRootNode(ctx, sqlc.CreateRootNodeParams{
			ID:        id,
			Name:      n.Name,
			StatusID:  n.StatusID,
			CreatedBy: n.CreatedBy,
		})
	} else {
		node, err = q.AppendNode(ctx, sqlc.AppendNodeParams{
			ID:        id,
			Type:      n.Type,
			Name:      n.Name,
			StatusID:  n.StatusID,
			CreatedBy: n.CreatedBy,
			ParentID:  *n.ParentID,
		})
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.OrgNode{}, appError.New(appError.KindNotFound, fmt.Sprintf("Parent not found or cannot contain a %s", n.Type))
		}
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return sqlc.OrgNode{}, appError.New(appError.KindConflict, fmt.Sprintf("A %s with this name already exists", n.Type))
	}
	if err != nil {
		return sqlc.OrgNode{}, fmt.Errorf("create org node: %w", err)
	}

	return node, nil
}
