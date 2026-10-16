package refresh

import (
	"context"

	"github.com/caneel7/keynote-api/internal/ent/gen"
	"github.com/caneel7/keynote-api/internal/features/auth/sessions"
)

type Service struct {
	db       *gen.Client
	sessions *sessions.Service
}

func NewService(db *gen.Client, sessions *sessions.Service) *Service {
	return &Service{db: db, sessions: sessions}
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (sessions.Tokens, error) {
	return s.sessions.Rotate(ctx, s.db, refreshToken)
}
