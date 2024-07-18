package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/internal/models"
	"instagram-clone.com/m/internal/session"
)

const (
	basePrefix = "session:"
)

type sessionRepository struct {
	redisClient *redis.Client
	basePrefix  string
	cfg         *config.Config
}

func NewSessionRepository(redisClient *redis.Client, cfg *config.Config) session.SessionRepository {
	return &sessionRepository{
		redisClient: redisClient,
		basePrefix:  basePrefix,
		cfg:         cfg,
	}
}

func (s *sessionRepository) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	session.SessionID = uuid.New().String()
	sessionKey := s.createKey(session.SessionID)

	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		return "", err
	}

	if err := s.redisClient.Set(ctx, sessionKey, sessionBytes, time.Duration(expire)*time.Second).Err(); err != nil {
		return "", err
	}

	return sessionKey, nil
}

func (s *sessionRepository) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	sessionBytes, err := s.redisClient.Get(ctx, sessionID).Bytes()
	if err != nil {
		return nil, err
	}

	session := &models.Session{}
	if err := json.Unmarshal(sessionBytes, &session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionRepository) DeleteSessionByID(ctx context.Context, sessionID string) error {
	if err := s.redisClient.Del(ctx, sessionID).Err(); err != nil {
		return err
	}
	return nil
}

func (s *sessionRepository) createKey(sessionID string) string {
	return fmt.Sprintf("%s:%s", s.basePrefix, sessionID)
}
