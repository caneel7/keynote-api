package sessions

import (
	"context"
	"fmt"
	"time"

	"github.com/caneel7/keynote-api/internal/common/tokens"
	"github.com/caneel7/keynote-api/internal/ent/gen"
	appError "github.com/caneel7/keynote-api/internal/errors"
	"github.com/google/uuid"
)

type Service struct {
	db *gen.Client
	signer *tokens.Signer
}

type TokensAndSession struct {
	AccessToken          string
	RefreshToken         string
	ExpireAT             time.Time
	RefreshTokenExpireAt time.Duration
}

func NewService(db *gen.Client,signer *tokens.Signer) *Service {
	return &Service{db,signer}
}

func (s *Service) CreateSessionAndTokens(
	ctx context.Context,
	userID uuid.UUID,
) (TokensAndSession, error) {

	var output TokensAndSession

	refresh, refreshHash, refreshttl, err := tokens.NewRefreshToken()
	if err != nil {
		return output, appError.InvalidCredentials()
	}

	session, err := s.db.UserSession.Create().
		SetUserID(userID).
		SetRefreshTokenHash(refreshHash).
		SetExpiresAt(time.Now().Add(time.Duration(refreshttl))).
		Save(ctx)

	if err != nil {
		return output, appError.InvalidCredentials()
	}

	token, exp, err := s.signer.Sign(userID, session.ID)
	if err != nil {
		return output, fmt.Errorf("Login error : %v", err)
	}

	output.ExpireAT = exp
	output.AccessToken = token
	output.RefreshToken = refresh
	output.RefreshTokenExpireAt = refreshttl

	return output, nil
}