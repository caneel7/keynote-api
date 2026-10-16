package refresh

import (
	"context"
	"time"

	"github.com/caneel7/keynote-api/internal/common/tokens"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/ent/gen/usersession"
	"github.com/caneel7/keynote-api/internal/errors"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
)

type Service struct {
	db     *gen.Client
	signer *tokens.Signer
}

func NewService(db *gen.Client, signer *tokens.Signer) *Service {
	return &Service{db, signer}
}

func (s *Service) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (sessions.TokensAndSession, error) {

	var output sessions.TokensAndSession

	sessionWithUser, err := s.db.UserSession.Query().
		Where(
			usersession.RefreshTokenHash(tokens.Hash(refreshToken)),
			usersession.RevokedAtIsNil(),
			usersession.ExpiresAtGT(time.Now()),
		).
		WithUser().
		Only(ctx)

	if err != nil {
		if gen.IsNotFound(err) {
			return output, errors.InvalidCredentials()
		}
		return output, errors.InvalidCredentials()
	}

	user := sessionWithUser.Edges.User

	refresh, refreshHash, refreshttl, err := tokens.NewRefreshToken()
	if err != nil {
		return output, errors.InvalidCredentials()
	}

	_, err = s.db.UserSession.UpdateOneID(sessionWithUser.ID).
		SetRefreshTokenHash(refreshHash).
		SetExpiresAt(time.Now().Add(time.Duration(refreshttl))).
		Save(ctx)

	if err != nil {
		return output, errors.InvalidCredentials()
	}

	token, exp, err := s.signer.Sign(user.ID, sessionWithUser.ID)
	if err != nil {
		return output, errors.InvalidCredentials()
	}

	output.AccessToken = token
	output.RefreshToken = refresh
	output.ExpireAT = exp
	output.RefreshTokenExpireAt = refreshttl

	return output, nil
}
