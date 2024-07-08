package repositories

import "instagram-clone.com/m/user/entities"

type UserRepository interface {
	GetUser(password string) (entities.User, error)
	CreateNewUser(input *entities.User) (entities.User, error)
	IsDuplicatedEmail(email string) bool
}
