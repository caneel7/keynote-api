package register

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	argon "github.com/caneel7/keynote-api/internal/common/password"
	"github.com/caneel7/keynote-api/internal/common/tokens"
	"github.com/caneel7/keynote-api/internal/constants"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/errors"
	appError "github.com/caneel7/keynote-api/internal/errors"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
)

type Service struct {
	db     *gen.Client
	signer *tokens.Signer
}

func NewService(db *gen.Client, signer *tokens.Signer) *Service {
	return &Service{db, signer}
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
) (sessions.TokensAndSession, error) {

	var output sessions.TokensAndSession

	email = strings.ToLower(strings.TrimSpace(email))

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return output, appError.Invalid(map[string]string{"email": "invalid email"})
	}

	hash, err := argon.HashPassword(password)
	if err != nil {
		return output, err
	}

	tx, err := s.db.Tx(ctx)
	if err != nil {
		return output, errors.New(appError.KindInternal, err.Error())
	}
	defer tx.Rollback()

	user, err := tx.User.Create().
		SetEmail(email).
		SetStateID(constants.StateIDs[constants.KindUsers][constants.StateActive]).
		Save(ctx)

	if err != nil {
		if gen.IsConstraintError(err) {
			return output, appError.New(
				appError.KindConflict,
				"Email Already Exists",
			)
		}
		return output, fmt.Errorf("Error creating user : %v", err)
	}

	if _, err := tx.UserCredential.Create().SetUsers(user).SetPasswordHash(hash).Save(ctx); err != nil {
		return output, fmt.Errorf("Error saving user credentials : %v", err)
	}

	if err := tx.Commit(); err != nil {
		return output, fmt.Errorf("Error commiting to db : %v", err)
	}

	out, err := sessions.NewService(s.db,s.signer).CreateSessionAndTokens(ctx,user.ID)
	if err != nil {
		return output, err
	}

	return out, nil
}
