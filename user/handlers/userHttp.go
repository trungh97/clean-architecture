package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"instagram-clone.com/m/user/usecases"
)

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type userHttpHandler struct {
	userUsecase usecases.UserUsecase
}

func NewUserHttpHandler(userUsecase usecases.UserUsecase) UserHandler {
	return &userHttpHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHttpHandler) Login(c echo.Context) error {
	reqBody := new(SignInInput)

	if err := c.Bind(reqBody); err != nil {
		log.Errorf("failed to bind request body: %v", err)
		return response(c, http.StatusBadRequest, "Bad Request", nil)
	}

	token, err := h.userUsecase.Login(reqBody.Email, reqBody.Password)
	if err != nil {
		log.Errorf("failed to login user: %v", err)
		return response(c, http.StatusInternalServerError, err.Error(), nil)
	}

	res := &SignInResponse{
		Token: token,
	}

	return response(c, http.StatusOK, "Login Successfully", res)
}
