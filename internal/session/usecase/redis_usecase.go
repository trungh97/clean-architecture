package usecase

import (
	"context"

	"instagram-clone.com/m/config"
	"instagram-clone.com/m/internal/models"
	"instagram-clone.com/m/internal/session"
)

type sessionUsecase struct {
	sessionRepo session.SessionRepository
	cfg         *config.Config
}

func NewSessionUsecase(sessionRepo session.SessionRepository, cfg *config.Config) session.SessionUsecase {
	return &sessionUsecase{
		sessionRepo: sessionRepo,
		cfg:         cfg,
	}
}

func (s *sessionUsecase) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	return s.sessionRepo.CreateSession(ctx, session, expire)
}

func (s *sessionUsecase) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	return s.sessionRepo.GetSessionByID(ctx, sessionID)
}

func (s *sessionUsecase) DeleteSessionByID(ctx context.Context, sessionID string) error {
	return s.sessionRepo.DeleteSessionByID(ctx, sessionID)
}
