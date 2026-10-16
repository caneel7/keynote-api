package register

import (
	"context"
	"fmt"
	"net/mail"
	"strings"

	argon "github.com/caneel7/keynote-api/internal/common/password"
	"github.com/caneel7/keynote-api/internal/constants"
	"github.com/caneel7/keynote-api/internal/ent/gen"
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

func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
) (sessions.Tokens, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return sessions.Tokens{}, appError.Invalid(map[string]string{"email": "invalid email"})
	}

	hash, err := argon.HashPassword(password)
	if err != nil {
		return sessions.Tokens{}, fmt.Errorf("hash password: %w", err)
	}

	tx, err := s.db.Tx(ctx)
	if err != nil {
		return sessions.Tokens{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	user, err := tx.User.Create().
		SetEmail(email).
		SetStateID(constants.StateIDs[constants.KindUsers][constants.StateActive]).
		Save(ctx)
	if gen.IsConstraintError(err) {
		return sessions.Tokens{}, appError.New(appError.KindConflict, "Email Already Exists")
	}
	if err != nil {
		return sessions.Tokens{}, fmt.Errorf("create user: %w", err)
	}

	if _, err := tx.UserCredential.Create().SetUsers(user).SetPasswordHash(hash).Save(ctx); err != nil {
		return sessions.Tokens{}, fmt.Errorf("create user credential: %w", err)
	}

	t, err := s.sessions.Issue(ctx, tx.Client(), user.ID)
	if err != nil {
		return sessions.Tokens{}, err
	}

	if err := tx.Commit(); err != nil {
		return sessions.Tokens{}, fmt.Errorf("commit: %w", err)
	}

	return t, nil
}
