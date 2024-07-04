package handlers

import (
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	Sample(c echo.Context) error
	Login(c echo.Context) error
}
