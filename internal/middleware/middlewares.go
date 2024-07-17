package middleware

import (
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/internal/auth"
	"instagram-clone.com/m/internal/session"
	"instagram-clone.com/m/pkg/logger"
)

type MiddlewareManager struct {
	sessionUsecase session.SessionUsecase
	authUsecase    auth.AuthUsecase
	cfg            *config.Config
	logger         logger.Logger
	origins        []string
}

func NewMiddlewareManager(cfg *config.Config, logger logger.Logger, origin []string, sessionUsecase session.SessionUsecase, authUsecase auth.AuthUsecase) *MiddlewareManager {
	return &MiddlewareManager{
		sessionUsecase: sessionUsecase,
		authUsecase:    authUsecase,
		cfg:            cfg,
		origins:        origin,
		logger:         logger,
	}
}
