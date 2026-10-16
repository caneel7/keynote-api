package login

import (
	"context"
	"fmt"
	"strings"

	argon "github.com/caneel7/keynote-api/internal/common/password"
	"github.com/caneel7/keynote-api/internal/common/tokens"
	"github.com/caneel7/keynote-api/internal/constants"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/ent/gen/user"
	"github.com/caneel7/keynote-api/internal/ent/gen/usercredential"
	appError "github.com/caneel7/keynote-api/internal/errors"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
)

type Service struct {
	db     *gen.Client
	signer *tokens.Signer
}

func NewService(db *gen.Client, signer *tokens.Signer) *Service {
	return &Service{db: db, signer: signer}
}


var dummyHash = "dum"

func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
) (sessions.TokensAndSession, error) {

	var output sessions.TokensAndSession
	email = strings.ToLower(strings.TrimSpace(email))

	u, err := s.db.UserCredential.Query().
		Where(usercredential.HasUsersWith(
			user.Email(email),
			user.DeletedAtIsNil(),
			user.StatusIDIn(
				constants.StateIDs[constants.KindUsers][constants.StateActive],
			),
		)).
		WithUsers().
		Only(ctx)

	if err != nil && !gen.IsNotFound(err) {
		return output, appError.InvalidCredentials()
	}

	if u == nil {
		argon.VerifyPassword(password, dummyHash) // equalize timing
		return output, appError.InvalidCredentials()
	}

	user := u.Edges.Users

	ok := argon.VerifyPassword(password, u.PasswordHash)
	if !ok {
		return output, appError.InvalidCredentials()
	}

	out, err := sessions.NewService(s.db, s.signer).CreateSessionAndTokens(ctx, user.ID)
	if err != nil {
		return output, fmt.Errorf("Login error : %v", err)
	}

	return out, nil
}


