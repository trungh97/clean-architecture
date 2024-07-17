package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id" redis:"id" validate:"omitempty"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	FirstName string    `json:"first_name" db:"first_name" validate:"required,lte=32"`
	LastName  string    `json:"last_name" db:"last_name" validate:"required,lte=32"`
	Password  string    `json:"password,omitempty" db:"password" validate:"omitempty,required,gte=8"`
	CreatedAt time.Time `json:"created_at,omitempty" redis:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" redis:"updated_at" db:"updated_at"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)
	u.ID = uuid.New()

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

func (u *User) SantinizePassword() {
	u.Password = ""
}

func (u *User) PrepareUpdate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	return nil
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
