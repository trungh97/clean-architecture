package handlers

import (
	"github.com/labstack/echo/v4"
	"instagram-clone.com/m/internal/auth"
	"instagram-clone.com/m/internal/middleware"
)

func MapAuthRoutes(authGroup *echo.Group, h auth.Handlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/logout", h.Logout())
	authGroup.GET("/:user_id", h.GetUserByID())
	// authGroup.Use(middleware.AuthJWTMiddleware(authUC, cfg))
	// You can use above middleware as the way of auth using JWT token
	authGroup.Use(mw.AuthSessionMiddleware)
	authGroup.GET("/me", h.GetMe())
	authGroup.GET("/token", h.GetCSRFToken())
	authGroup.PUT("/:user_id", h.UpdateUser(), mw.CSRF)
	authGroup.DELETE("/:user_id", h.DeleteUser(), mw.CSRF)
}
