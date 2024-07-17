package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"instagram-clone.com/m/database"
	"instagram-clone.com/m/internal/auth"
	"instagram-clone.com/m/internal/models"
)

type AuthRepository struct {
	db database.Database
}

func NewAuthRepository(db database.Database) auth.AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}

	if err := r.db.GetDb().Where("id = ?", userID).First(user).Error; err != nil {
		return nil, errors.New("failed to get user from the database")
	}

	return user, nil
}

func (r *AuthRepository) IsDuplicatedEmail(email string) bool {
	user := &models.User{}

	err := r.db.GetDb().Where("email = ?", email).First(&user).Error

	return err == nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := r.db.GetDb().Create(&models.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		ID:        user.ID,
	}).Error

	if err != nil {
		formattedError := fmt.Sprintf("failed to create user in the database: %v", err)
		log.Errorf("failed to create user in the database: %v", err)
		return nil, errors.New(formattedError)
	}

	return user, err
}

func (r *AuthRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := r.db.GetDb().Model(&models.User{}).Where("id = ?", user.ID).Updates(&models.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}).Error

	if err != nil {
		formattedError := fmt.Sprintf("failed to update user in the database: %v", err)
		log.Errorf("failed to update user in the database: %v", err)
		return nil, errors.New(formattedError)
	}

	return user, nil
}

func (r *AuthRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {

	err := r.db.GetDb().Where("id = ?", userID).Delete(&models.User{}).Error

	if err != nil {
		formattedError := fmt.Sprintf("failed to delete user in the database: %v", err)
		log.Errorf("failed to delete user in the database: %v", err)
		return errors.New(formattedError)
	}

	return nil
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	foundUser := &models.User{}

	err := r.db.GetDb().Where("email = ?", user.Email).First(&foundUser).Error
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}
