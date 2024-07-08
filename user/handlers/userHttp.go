package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"instagram-clone.com/m/user/models"
	"instagram-clone.com/m/user/usecases"
)

type SignInResponse struct {
	Token string `json:"token"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
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
	reqBody := new(models.AuthInput)

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

func (h *userHttpHandler) Register(c echo.Context) error {
	reqBody := new(models.AuthInput)

	if err := c.Bind(reqBody); err != nil {
		log.Errorf("failed to bind request body: %v", err)
		return response(c, http.StatusBadRequest, "Bad Request", nil)
	}

	userId, err := h.userUsecase.Register(reqBody.Email, reqBody.Password)
	if err != nil {
		log.Errorf("failed to register user: %v", err)
		return response(c, http.StatusInternalServerError, err.Error(), nil)
	}

	res := &RegisterResponse{
		UserID: userId,
	}

	return response(c, http.StatusOK, "Register Successfully", res)
}
