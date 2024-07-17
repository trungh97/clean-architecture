package auth

import (
	"context"

	"instagram-clone.com/m/internal/models"
)

type RedisRepository interface {
	GetUserByIDCtx(ctx context.Context, key string) (*models.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error
	DeleteUserCtx(ctx context.Context, key string) error
}
