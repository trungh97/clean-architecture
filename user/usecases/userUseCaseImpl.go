package usecases

import (
	"instagram-clone.com/m/user/entities"
	"instagram-clone.com/m/user/models"
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

func (u *userUsecaseImpl) UserDataProcessing(in *models.AddUserData) error {
	insertUserData := &entities.AddUserData{
		Email:    in.Email,
		Password: in.Password,
	}

	if err := u.userRepository.InsertUserData(insertUserData); err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) UserLogin(in *models.AddUserData) (string, error) {
	loginUserData := &entities.AddUserData{
		Email:    in.Email,
		Password: in.Password,
	}

	token, err := u.userRepository.Login(loginUserData)
	if err != nil {
		return "", err
	}

	return token, nil
}
