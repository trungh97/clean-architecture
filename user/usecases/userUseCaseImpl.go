package usecases

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/user/repositories"
)

type userUsecaseImpl struct {
	userRepository repositories.UserRepository
}

func NewUserUseCaseImpl(userRepository repositories.UserRepository) *userUsecaseImpl {
	return &userUsecaseImpl{
		userRepository: userRepository,
	}
}

func (u *userUsecaseImpl) Login(email, password string) (string, error) {
	user, err := u.userRepository.GetUser(email)

	if err != nil {
		// If the user is not found, return an error
		return "", err
	}

	// Check if the provided password is correct
	// Compare the hashed password in the database with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println(user.Password, password)
		// If the passwords do not match, return an error
		return "", errors.New("invalid email or password")
	}

	token, err := generateToken(user.Email)
	if err != nil {
		log.Errorf("failed to generate token: %v", err)
		return "", errors.New("internal server error")
	}

	return token, nil
}

func generateToken(email string) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"iss": "instagram-clone.com",
		"exp": 3600,
	})

	tokenString, err := claims.SignedString([]byte(config.GetConfig().JWT.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
