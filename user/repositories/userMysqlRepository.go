package repositories

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/database"
	"instagram-clone.com/m/user/entities"
)

type UserMysqlRepository struct {
	db database.Database
}

func NewUserMysqlRepository(db database.Database) UserRepository {
	return &UserMysqlRepository{
		db: db,
	}
}

func (r *UserMysqlRepository) InsertUserData(in *entities.AddUserData) error {
	data := &entities.User{
		Email:    in.Email,
		Password: in.Password,
	}

	result := r.db.GetDb().Create(data)
	if result.Error != nil {
		log.Errorf("failed to insert user data: %v", result.Error)
		return result.Error
	}

	log.Debugf("InsertUserData result: %v", result.RowsAffected)

	return nil
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

func (r *UserMysqlRepository) Login(in *entities.AddUserData) (string, error) {
	var user entities.User
	data := &entities.User{
		Email:    in.Email,
		Password: in.Password,
	}

	err := r.db.GetDb().Where("email = ?", data.Email).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("invalid username or password")
		}
		log.Errorf("failed to get user from the database: %v", err)
		return "", errors.New("internal server error")
	}

	// Check if the provided password is correct
	// Compare the hashed password in the database with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		fmt.Println(user.Password, data.Password)
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
