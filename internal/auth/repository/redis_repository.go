package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"instagram-clone.com/m/internal/auth"
	"instagram-clone.com/m/internal/models"
)

type authRedisRepository struct {
	redisClient *redis.Client
}

func NewAuthRedisRepository(redisClient *redis.Client) auth.RedisRepository {
	return &authRedisRepository{redisClient: redisClient}
}

func (r *authRedisRepository) GetUserByIDCtx(ctx context.Context, key string) (*models.User, error) {
	userBytes, err := r.redisClient.Get(ctx, key).Bytes()

	if err != nil {
		return nil, err
	}

	user := &models.User{}

	if err := json.Unmarshal(userBytes, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *authRedisRepository) SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) error {
	userBytes, err := json.Marshal(user)

	if err != nil {
		return err
	}

	if err := r.redisClient.Set(ctx, key, userBytes, time.Duration(seconds)*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (r *authRedisRepository) DeleteUserCtx(ctx context.Context, key string) error {
	if err := r.redisClient.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
