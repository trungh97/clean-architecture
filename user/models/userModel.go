package models

type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
