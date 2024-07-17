package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
	"instagram-clone.com/m/docs"
	"instagram-clone.com/m/pkg/csrf"
	"instagram-clone.com/m/pkg/utils"

	authHttp "instagram-clone.com/m/internal/auth/handlers"
	authRepository "instagram-clone.com/m/internal/auth/repository"
	authUsecase "instagram-clone.com/m/internal/auth/usecase"

	sessionRepository "instagram-clone.com/m/internal/session/repository"
	sessionUsecase "instagram-clone.com/m/internal/session/usecase"

	apiMiddlewares "instagram-clone.com/m/internal/middleware"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init repositories
	authRepo := authRepository.NewAuthRepository(s.db)
	authRedisRepo := authRepository.NewAuthRedisRepository(s.redisClient)
	sessionRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)

	// Init usecases
	authUC := authUsecase.NewAuthUsecase(authRepo, authRedisRepo, s.cfg, s.logger)
	sessionUC := sessionUsecase.NewSessionUsecase(sessionRepo, s.cfg)

	// Init handlers
	authHandler := authHttp.NewAuthHandler(s.cfg, authUC, sessionUC, s.logger)

	// Middlewares
	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, s.logger, []string{"*"}, sessionUC, authUC)

	docs.SwaggerInfo.Title = "Instagram Clone BE API"
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	v1 := e.Group("/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")

	authHttp.MapAuthRoutes(authGroup, authHandler, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check Request ID: %s", utils.GetRequestID(c))

		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
