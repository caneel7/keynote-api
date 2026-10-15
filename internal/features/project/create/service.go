package create

import (
	"context"

	"github.com/caneel7/keynote-api/internal/constants"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/ent/gen/project"
	"github.com/caneel7/keynote-api/internal/errors"
	"github.com/google/uuid"
)

type Service struct {
	db *gen.Client
}

func NewService(db *gen.Client) (*Service) {
	return &Service{db}
}

func (s *Service) CreateProject(
	ctx context.Context,
	name string,
	tenantID uuid.UUID,
	workspaceID uuid.UUID,
	userID uuid.UUID,
) (gen.Project,error) {

	exist, err := s.db.Project.Query().
		Where(
			project.Name(name),
			project.TenantID(tenantID),
			project.WorkspaceID(workspaceID),
		).
		Exist(ctx);

	if err != nil {
		return gen.Project{}, errors.New(errors.KindInternal, err.Error())
	}

	if exist {
		return gen.Project{},errors.New(errors.KindConflict,"Project with name exists")
	}

	project, err := s.db.Project.Create().
		SetName(name).
		SetTenantID(tenantID).
		SetWorkspaceID(workspaceID).
		SetStateID(constants.StateIDs[constants.KindProject][constants.StateActive]).
		SetCreatorID(userID).
		Save(ctx)

	if err != nil {
		if gen.IsConstraintError(err){
			return gen.Project{},errors.New(errors.KindConflict,err.Error())
		}
		return gen.Project{}, errors.New(errors.KindInternal, "Failed to create project")
	}

	return *project,nil
}