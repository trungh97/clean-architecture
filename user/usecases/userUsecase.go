package usecases

import "instagram-clone.com/m/user/models"

type UserUsecase interface {
	UserDataProcessing(in *models.AddUserData) error
	UserLogin(in *models.AddUserData) (string, error)
}
