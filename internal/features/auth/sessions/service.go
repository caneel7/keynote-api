package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/caneel7/keynote-api/internal/common/tokens"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/ent/gen/user"
	"github.com/caneel7/keynote-api/internal/ent/gen/usersession"
	appError "github.com/caneel7/keynote-api/internal/errors"
	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken      string
	AccessExpiresAt  time.Time
	RefreshToken     string
	RefreshExpiresAt time.Time
}

type Service struct {
	signer *tokens.Signer
}

func NewService(signer *tokens.Signer) *Service {
	return &Service{signer: signer}
}

func (s *Service) Issue(ctx context.Context, db *gen.Client, userID uuid.UUID) (Tokens, error) {
	refresh, refreshHash, ttl, err := tokens.NewRefreshToken()
	if err != nil {
		return Tokens{}, fmt.Errorf("generate refresh token: %w", err)
	}
	refreshExpiresAt := time.Now().Add(ttl)

	session, err := db.UserSession.Create().
		SetUserID(userID).
		SetRefreshTokenHash(refreshHash).
		SetExpiresAt(refreshExpiresAt).
		Save(ctx)
	if err != nil {
		return Tokens{}, fmt.Errorf("create session: %w", err)
	}

	access, accessExpiresAt, err := s.signer.Sign(userID, session.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("sign access token: %w", err)
	}

	return Tokens{
		AccessToken:      access,
		AccessExpiresAt:  accessExpiresAt,
		RefreshToken:     refresh,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func (s *Service) Rotate(ctx context.Context, db *gen.Client, refreshToken string) (Tokens, error) {
	now := time.Now()
	oldHash := tokens.Hash(refreshToken)

	session, err := db.UserSession.Query().
		Where(
			usersession.RefreshTokenHash(oldHash),
			usersession.RevokedAtIsNil(),
			usersession.ExpiresAtGT(now),
			usersession.HasUserWith(user.DeletedAtIsNil()),
		).
		Only(ctx)
	if gen.IsNotFound(err) {
		return Tokens{}, invalidRefreshToken()
	}
	if err != nil {
		return Tokens{}, fmt.Errorf("find session: %w", err)
	}

	refresh, refreshHash, ttl, err := tokens.NewRefreshToken()
	if err != nil {
		return Tokens{}, fmt.Errorf("generate refresh token: %w", err)
	}
	refreshExpiresAt := now.Add(ttl)

	updated, err := db.UserSession.Update().
		Where(
			usersession.ID(session.ID),
			usersession.RefreshTokenHash(oldHash),
			usersession.RevokedAtIsNil(),
		).
		SetRefreshTokenHash(refreshHash).
		SetExpiresAt(refreshExpiresAt).
		SetLastUsedAt(now).
		Save(ctx)
	if err != nil {
		return Tokens{}, fmt.Errorf("rotate session: %w", err)
	}
	if updated == 0 {
		return Tokens{}, invalidRefreshToken()
	}

	access, accessExpiresAt, err := s.signer.Sign(session.UserID, session.ID)
	if err != nil {
		return Tokens{}, fmt.Errorf("sign access token: %w", err)
	}

	return Tokens{
		AccessToken:      access,
		AccessExpiresAt:  accessExpiresAt,
		RefreshToken:     refresh,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func invalidRefreshToken() *appError.AppError {
	return appError.New(appError.KindUnauthenticated, "Invalid or expired refresh token")
}
