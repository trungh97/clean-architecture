package handlers

import (
	"strings"

	"github.com/labstack/echo/v4"
)

type BaseResponse struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

func BuildResponse(status bool, message string, data interface{}) BaseResponse {
	res := BaseResponse{
		Message: message,
		Errors:  nil,
		Data:    data,
	}
	return res
}

func response(c echo.Context, responseCode int, message string, data interface{}, err string) error {
	var splittedError []string
	if err == "" {
		splittedError = []string{}
	} else {
		splittedError = strings.Split(err, "\n")
	}
	return c.JSON(responseCode, &BaseResponse{
		Message: message,
		Errors:  splittedError,
		Data:    data,
	})
}
