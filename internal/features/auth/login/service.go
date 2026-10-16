package login

import (
	"context"
	"fmt"
	"strings"

	argon "github.com/caneel7/keynote-api/internal/common/password"
	"github.com/caneel7/keynote-api/internal/constants"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/ent/gen/user"
	"github.com/caneel7/keynote-api/internal/ent/gen/usercredential"
	appError "github.com/caneel7/keynote-api/internal/errors"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
)

type Service struct {
	db       *gen.Client
	sessions *sessions.Service
}

func NewService(db *gen.Client, sessions *sessions.Service) *Service {
	return &Service{db: db, sessions: sessions}
}

var dummyHash = mustHash("keynote-dummy-password")

func mustHash(password string) string {
	h, err := argon.HashPassword(password)
	if err != nil {
		panic(fmt.Sprintf("login: compute dummy hash: %v", err))
	}
	return h
}

func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
) (sessions.Tokens, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	cred, err := s.db.UserCredential.Query().
		Where(usercredential.HasUsersWith(
			user.Email(email),
			user.DeletedAtIsNil(),
			user.StatusID(constants.StateIDs[constants.KindUsers][constants.StateActive]),
		)).
		Only(ctx)
	if err != nil && !gen.IsNotFound(err) {
		return sessions.Tokens{}, fmt.Errorf("find credential: %w", err)
	}

	if cred == nil {
		argon.VerifyPassword(password, dummyHash) // equalize timing
		return sessions.Tokens{}, appError.InvalidCredentials()
	}

	if !argon.VerifyPassword(password, cred.PasswordHash) {
		return sessions.Tokens{}, appError.InvalidCredentials()
	}

	return s.sessions.Issue(ctx, s.db, cred.UserID)
}
