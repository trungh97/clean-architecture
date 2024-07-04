package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthService provides authentication services
type AuthService struct {
	DB       *gorm.DB
	FireAuth *auth.Client
}

// Login authenticates a user with the provided email and password
// and generates a Firebase custom token for the user.
//
// It returns the Firebase custom token and an error if any.
func (s *AuthService) Login(email, password string) (string, error) {
	// Get the user from the database
	var user User
	// Query the database for a user with the given email
	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		// If the user is not found, return an error
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("invalid email or password")
		}
		// If there is an error, log it and return an error
		log.Printf("failed to get user by email from the database: %v", err)
		return "", errors.New("internal server error")
	}

	// Check if the provided password is correct
	// Compare the hashed password in the database with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// If the passwords do not match, return an error
		return "", errors.New("invalid email or password")
	}

	// Generate a Firebase custom token for the user
	// Generate a custom token for the user with the given ID
	token, err := s.FireAuth.CustomToken(context.Background(), user.ID)
	if err != nil {
		// If there is an error, log it and return an error
		log.Printf("failed to generate custom token for user %s: %v", user.ID, err)
		return "", errors.New("internal server error")
	}

	return token, nil
}

func (s *AuthService) Register(email, password string) (string, error) {
	var user User

	if err := s.DB.Raw("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("failed to get user by email from the database: %v", err)
			return "", errors.New("internal server error")
		}
	}

	if user.ID != "" {
		return "", errors.New("user already exists")
	}

	// Generate a UUID for the new user
	uid := uuid.New().String()

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return "", errors.New("internal server error")
	}

	// Create a new user in the database
	if err := s.DB.Exec("INSERT INTO users (id, email, password) VALUES (?, ?, ?)", uid, email, hashedPassword).Error; err != nil {
		log.Printf("failed to insert user into database: %v", err)
		return "", errors.New("internal server error")
	}

	// Create a custom token for the user using the Firebase Admin SDK
	customToken, err := s.FireAuth.CustomToken(context.Background(), uid)
	if err != nil {
		log.Printf("failed to create custom token for user: %v", err)
		return "", errors.New("internal server error")
	}

	return customToken, nil
}

func generateJWT() (string, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {

}
