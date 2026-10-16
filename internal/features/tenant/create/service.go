package create

import (
	"context"
	"fmt"

	"github.com/caneel7/keynote-api/internal/common/slug"
	"github.com/caneel7/keynote-api/internal/constants"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/ent/gen/tenant"
	"github.com/caneel7/keynote-api/internal/errors"
	"github.com/google/uuid"
)

type Service struct {
	db *gen.Client
}

func NewService(db *gen.Client) (*Service) {
	return &Service{db}
}

func (s *Service) Create(
	ctx context.Context,
	name string,
	userID uuid.UUID,
) (gen.Tenant,gen.Workspace,error) {

	var newTenant gen.Tenant
	var newWorkspace gen.Workspace

	existingTenant, err := s.db.Tenant.Query().
		Where(tenant.CreatedBy(userID)).
		Exist(ctx)

	if err != nil {
		return newTenant,newWorkspace, fmt.Errorf("Error executing query %v",err)
	}

	if existingTenant {
		return newTenant,newWorkspace, errors.New(errors.KindBadRequest,"Cannot create more than one tenant")
	}

	slug, err := s.createTenantSlug(ctx, name)
	if err != nil {
		return newTenant,newWorkspace, err
	}

	tx, err := s.db.Tx(ctx)
	if err != nil {
		return newTenant,newWorkspace, err
	}
	defer tx.Rollback()

	newtenant, err := tx.Tenant.Create().
		SetName(name).
		SetSlug(slug).
		SetCreatedBy(userID).
		SetStateID(constants.StateIDs[constants.KindTenant][constants.StateActive]).
		Save(ctx)

	if err != nil {
		if gen.IsConstraintError(err) {
			return newTenant,newWorkspace,errors.New(errors.KindConflict,err.Error())
		}
		_ = fmt.Errorf("Error Creating project : %v",err.Error())
		return newTenant,newWorkspace, errors.New(errors.KindInternal,"Error creating project. Please start again later")
	}

	workspace, err := tx.Workspace.Create().
		SetName(name).
		SetSlug(slug).
		SetTenantID(newtenant.ID).
		SetCreatedBy(userID).
		SetStatusID(constants.StateIDs[constants.KindWorkspace][constants.StateActive]).
		Save(ctx)

	if err != nil {
		if gen.IsConstraintError(err) {
			return newTenant,newWorkspace,errors.New(errors.KindConflict,err.Error())
		}
		return newTenant,newWorkspace, err
	}
	
	if err := tx.Commit(); err != nil {
		return newTenant,newWorkspace, fmt.Errorf("Error commiting to db : %v", err)
	}

	return *newtenant,*workspace,nil

}

func (s *Service) createTenantSlug(ctx context.Context,name string) (string,error) {
	
	base := slug.Slugify(name)

	for i := 0; i < 5; i++ {
		exist, err := s.db.Tenant.Query().
			Where(tenant.Slug(base)).Exist(ctx)
		
		if err != nil {
			return "",err
		}

		if !exist {
			return base,nil
		}
		base += "-" + slug.Randomsuffix()
	}
	return "",errors.New(errors.KindConflict,"Cannot find free slug")
}