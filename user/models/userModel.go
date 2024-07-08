package models

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}
