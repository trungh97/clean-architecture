package usecases

type UserUsecase interface {
	Login(email, password string) (string, error)
}
