package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"instagram-clone.com/m/user/models"
	"instagram-clone.com/m/user/usecases"
)

type userHttpHandler struct {
	userUsecase usecases.UserUsecase
}

func NewUserHttpHandler(userUsecase usecases.UserUsecase) UserHandler {
	return &userHttpHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHttpHandler) Sample(c echo.Context) error {
	reqBody := new(models.AddUserData)

	if err := c.Bind(reqBody); err != nil {
		log.Errorf("failed to bind request body: %v", err)
		return response(c, http.StatusBadRequest, "Bad Request", nil)
	}

	if err := h.userUsecase.UserDataProcessing(reqBody); err != nil {
		log.Errorf("failed to insert user data: %v", err)
		return response(c, http.StatusInternalServerError, "Processing Data Failed", nil)
	}

	return response(c, http.StatusOK, "Success", nil)
}

func (h *userHttpHandler) Login(c echo.Context) error {
	reqBody := new(models.AddUserData)

	if err := c.Bind(reqBody); err != nil {
		log.Errorf("failed to bind request body: %v", err)
		return response(c, http.StatusBadRequest, "Bad Request", nil)
	}

	token, err := h.userUsecase.UserLogin(reqBody)
	if err != nil {
		log.Errorf("failed to login user: %v", err)
		return response(c, http.StatusInternalServerError, err.Error(), nil)
	}

	res := &LoginResponse{
		Token: token,
	}

	return response(c, http.StatusOK, "Login Successfully", res)
}
