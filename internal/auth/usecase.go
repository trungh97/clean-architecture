package auth

import (
	"context"

	"github.com/google/uuid"
	"instagram-clone.com/m/internal/models"
)

type AuthUsecase interface {
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Register(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUserByID(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
