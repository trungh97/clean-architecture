package repositories

import "instagram-clone.com/m/user/entities"

type UserRepository interface {
	InsertUserData(in *entities.AddUserData) error
	Login(in *entities.AddUserData) (string, error)
}
