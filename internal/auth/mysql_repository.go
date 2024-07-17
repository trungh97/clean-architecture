package auth

import (
	"context"

	"github.com/google/uuid"
	"instagram-clone.com/m/internal/models"
)

type AuthRepository interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	FindUserByEmail(ctx context.Context, user *models.User) (*models.User, error)
	IsDuplicatedEmail(email string) bool
}
