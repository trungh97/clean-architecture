package usecases

type UserUsecase interface {
	Login(email, password string) (string, error)
	Register(email, username, password string) (string, error)
}
