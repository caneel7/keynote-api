package orgnodes

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/ent/gen/orgnode"
	"github.com/caneel7/keynote-api/internal/ent/gen/predicate"
	appError "github.com/caneel7/keynote-api/internal/errors"
	"github.com/google/uuid"
)

type NewNode struct {
	ParentID  *uuid.UUID
	Type      orgnode.Type
	Name      string
	StatusID  uuid.UUID
	CreatedBy uuid.UUID
}

var childType = map[orgnode.Type]orgnode.Type{
	orgnode.TypeTenant:    orgnode.TypeWorkspace,
	orgnode.TypeWorkspace: orgnode.TypeProject,
}

func AppendNode(ctx context.Context, db *gen.Client, n NewNode) (*gen.OrgNode, error) {
	id := uuid.New()
	create := db.OrgNode.Create().
		SetID(id).
		SetType(n.Type).
		SetName(n.Name).
		SetStatusID(n.StatusID).
		SetCreatedBy(n.CreatedBy)

	if n.ParentID == nil {
		if n.Type != orgnode.TypeTenant {
			return nil, appError.New(appError.KindBadRequest, "Only a tenant can be created without a parent")
		}
		create.SetTenantID(id).SetPath(label(id))
	} else {
		parent, err := db.OrgNode.Query().
			Where(orgnode.ID(*n.ParentID), orgnode.DeletedAtIsNil()).
			Only(ctx)
		if gen.IsNotFound(err) {
			return nil, appError.New(appError.KindNotFound, "Parent not found")
		}
		if err != nil {
			return nil, fmt.Errorf("find parent node: %w", err)
		}
		if childType[parent.Type] != n.Type {
			return nil, appError.New(appError.KindBadRequest, fmt.Sprintf("A %s cannot be created under a %s", n.Type, parent.Type))
		}
		create.SetTenantID(parent.TenantID).
			SetParentID(parent.ID).
			SetPath(parent.Path + "." + label(id))
	}

	node, err := create.Save(ctx)
	if gen.IsConstraintError(err) {
		return nil, appError.New(appError.KindConflict, fmt.Sprintf("A %s with this name already exists", n.Type))
	}
	if err != nil {
		return nil, fmt.Errorf("create org node: %w", err)
	}
	return node, nil
}

func label(id uuid.UUID) string {
	return id.String()
}

func IsAncestorOf(path string) predicate.OrgNode {
	return func(s *sql.Selector) {
		s.Where(sql.P(func(b *sql.Builder) {
			b.Ident(s.C(orgnode.FieldPath)).WriteString(" @> ").Arg(path).WriteString("::ltree")
		}))
	}
}

func IsDescendantOf(path string) predicate.OrgNode {
	return func(s *sql.Selector) {
		s.Where(sql.P(func(b *sql.Builder) {
			b.Ident(s.C(orgnode.FieldPath)).WriteString(" <@ ").Arg(path).WriteString("::ltree")
		}))
	}
}
