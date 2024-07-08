package server

import (
	"fmt"

	userHandlers "instagram-clone.com/m/user/handlers"
	userRepositories "instagram-clone.com/m/user/repositories"
	userUseCases "instagram-clone.com/m/user/usecases"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/database"
)

type echoServer struct {
	app  *echo.Echo
	db   database.Database
	conf *config.Config
}

type (
	CustomValidator struct {
		validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

func NewEchoServer(conf *config.Config, db database.Database) *echoServer {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG)
	echoApp.Validator = &CustomValidator{validator: validator.New()}

	return &echoServer{
		app:  echoApp,
		db:   db,
		conf: conf,
	}
}

func (s *echoServer) Start() {
	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())

	// Health check adding
	s.app.GET("v1/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	s.initializeUserHttpHandler()

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	s.app.Logger.Fatal(s.app.Start(serverUrl))
}

func (s *echoServer) initializeUserHttpHandler() {
	// Initialize all layers
	userMysqlRepository := userRepositories.NewUserMysqlRepository(s.db)

	userUseCase := userUseCases.NewUserUseCaseImpl(userMysqlRepository)
	userHttpHandler := userHandlers.NewUserHttpHandler(userUseCase)

	// Router
	authRouters := s.app.Group("/v1/auth")

	authRouters.POST("/login", userHttpHandler.Login)
	authRouters.POST("/register", userHttpHandler.Register)
}
