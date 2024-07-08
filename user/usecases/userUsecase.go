package usecases

type UserUsecase interface {
	Login(email, password string) (string, error)
	Register(email, password string) (string, error)
}
