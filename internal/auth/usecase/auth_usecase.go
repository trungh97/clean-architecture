package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/internal/auth"
	"instagram-clone.com/m/internal/models"
	"instagram-clone.com/m/pkg/httpErrors"
	"instagram-clone.com/m/pkg/logger"
	"instagram-clone.com/m/pkg/utils"
)

type authUsecase struct {
	authRepository  auth.AuthRepository
	redisRepository auth.RedisRepository
	cfg             *config.Config
	logger          logger.Logger
}

const (
	basePrefix    = "api-auth:"
	cacheDuration = 3600
)

func NewAuthUsecase(authRepository auth.AuthRepository, redisRepository auth.RedisRepository, cfg *config.Config, logger logger.Logger) auth.AuthUsecase {
	return &authUsecase{
		authRepository:  authRepository,
		redisRepository: redisRepository,
		cfg:             cfg,
		logger:          logger,
	}
}

func (u *authUsecase) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	foundUser, err := u.authRepository.FindUserByEmail(ctx, user)

	if err != nil {
		// If the user is not found, return an error
		return nil, err
	}

	// Check if the provided password is correct
	// Compare the hashed password in the database with the provided password
	if err := foundUser.ComparePassword(user.Password); err != nil {
		// If the passwords do not match, return an error
		return nil, httpErrors.NewUnauthorizedError(err)
	}

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(err)
	}

	foundUser.SantinizePassword()

	userWithToken := &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}

	return userWithToken, nil
}

func (u *authUsecase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	existedUser, err := u.authRepository.FindUserByEmail(ctx, user)
	if existedUser != nil || err == nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrEmailAlreadyExists, nil)
	}

	if err := user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	createdUser, err := u.authRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser.SantinizePassword()

	return createdUser, nil
}

func (u *authUsecase) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := user.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(err)
	}

	updatedUser, err := u.authRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser.SantinizePassword()

	if err := u.redisRepository.SetUserCtx(ctx, u.GenerateUserKey(updatedUser.ID.String()), cacheDuration, updatedUser); err != nil {
		u.logger.Errorf("AuthUsecase.Update.SetUserCtx: %s", err)
	}

	return updatedUser, nil
}

func (u *authUsecase) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	if err := u.authRepository.DeleteUser(ctx, userID); err != nil {
		return err
	}

	if err := u.redisRepository.DeleteUserCtx(ctx, u.GenerateUserKey(userID.String())); err != nil {
		u.logger.Errorf("AuthUsecase.Delete.DeleteUserCtx: %s", err)
	}
	return nil
}

func (u *authUsecase) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	userKey := u.GenerateUserKey(userID.String())

	cachedUser, err := u.redisRepository.GetUserByIDCtx(ctx, userKey) // Get user from cache

	if err != nil {
		u.logger.Errorf("AuthUsecase.GetUserByID.GetUserByIDCtx: %s", err)
	}

	if cachedUser != nil {
		cachedUser.SantinizePassword()
		return cachedUser, nil
	}

	user, err := u.authRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.SantinizePassword()

	if err := u.redisRepository.SetUserCtx(ctx, userKey, cacheDuration, user); err != nil {
		u.logger.Errorf("AuthUsecase.GetUserByID.SetUserCtx: %s", err)
	}

	return user, nil
}

func (u *authUsecase) GenerateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}
