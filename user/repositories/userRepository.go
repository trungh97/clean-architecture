package repositories

import "instagram-clone.com/m/user/entities"

type UserRepository interface {
	GetUser(password string) (entities.User, error)
}
