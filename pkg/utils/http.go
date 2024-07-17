package utils

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/pkg/httpErrors"
	"instagram-clone.com/m/pkg/logger"
)

type UserCtxKey struct{}

// Get request id from echo context
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}

// Get user ip address
func GetIPAddress(c echo.Context) string {
	return c.Request().RemoteAddr
}

// Read request body and validate
func ReadRequest(ctx echo.Context, request interface{}) error {
	if err := ctx.Bind(request); err != nil {
		return err
	}
	return validate.StructCtx(ctx.Request().Context(), request)
}

// Error response with logging error for echo context
func LogResponseError(ctx echo.Context, err error, logger logger.Logger) {
	logger.Errorf("ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %v", GetRequestID(ctx), GetIPAddress(ctx), err)
}

// Configure jwt cookie
func CreateSessionCookie(cfg *config.Config, session string) *http.Cookie {
	return &http.Cookie{
		Name:  cfg.Session.Name,
		Value: session,
		Path:  "/",
		// Domain: "/",
		// Expires:    time.Now().Add(1 * time.Minute),
		RawExpires: "",
		MaxAge:     cfg.Session.Expire,
		Secure:     cfg.Cookie.Secure,
		HttpOnly:   cfg.Cookie.HTTPOnly,
		SameSite:   0,
	}
}

// Delete session
func DeleteSessionCookie(c echo.Context, sessionName string) {
	c.SetCookie(&http.Cookie{
		Name:   sessionName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

// Error response with logging error for echo context
func ErrResponseWithLog(ctx echo.Context, err error) error {
	log.Printf(
		"ErrResponseWithLog, RequestID: %s, IPAddress: %s, Error: %s",
		GetRequestID(ctx),
		GetIPAddress(ctx),
		err,
	)
	return ctx.JSON(httpErrors.ErrorResponse(err))
}
