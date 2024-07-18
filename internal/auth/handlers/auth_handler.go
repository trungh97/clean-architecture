package handlers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"instagram-clone.com/m/config"
	"instagram-clone.com/m/internal/auth"
	"instagram-clone.com/m/internal/models"
	"instagram-clone.com/m/internal/session"
	"instagram-clone.com/m/pkg/csrf"
	"instagram-clone.com/m/pkg/httpErrors"
	"instagram-clone.com/m/pkg/logger"
	"instagram-clone.com/m/pkg/utils"
)

type authHandler struct {
	cfg            *config.Config
	authUsecase    auth.AuthUsecase
	sessionUsecase session.SessionUsecase
	logger         logger.Logger
}

func NewAuthHandler(cfg *config.Config, authUsecase auth.AuthUsecase, sessionUsecase session.SessionUsecase, logger logger.Logger) auth.Handlers {
	return &authHandler{
		cfg:            cfg,
		authUsecase:    authUsecase,
		sessionUsecase: sessionUsecase,
		logger:         logger,
	}
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /auth/register [post]
func (h *authHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := &models.User{}

		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		createdUser, err := h.authUsecase.Register(c.Request().Context(), user)
		if err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		session, err := h.sessionUsecase.CreateSession(c.Request().Context(), &models.Session{
			UserID: createdUser.ID,
		}, h.cfg.Session.Expire)

		if err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		c.SetCookie(utils.CreateSessionCookie(h.cfg, session))

		return c.JSON(http.StatusCreated, createdUser)
	}
}

// Login godoc
// @Summary Login user
// @Description login user, returns user with token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/login [post]
func (h *authHandler) Login() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"required,email"`
		Password string `json:"password" db:"password" validate:"required,gte=8"`
	}

	return func(c echo.Context) error {
		reqBody := &Login{}

		if err := utils.ReadRequest(c, reqBody); err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		userWithToken, err := h.authUsecase.Login(c.Request().Context(), &models.User{
			Email:    reqBody.Email,
			Password: reqBody.Password,
		})

		if err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		session, err := h.sessionUsecase.CreateSession(c.Request().Context(), &models.Session{
			UserID: userWithToken.User.ID,
		}, h.cfg.Session.Expire)

		if err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		c.SetCookie(utils.CreateSessionCookie(h.cfg, session))

		return c.JSON(http.StatusOK, userWithToken)
	}
}

// Logout godoc
// @Summary Logout user
// @Description logout user removing session
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {string} string	"ok"
// @Router /auth/logout [post]
func (h *authHandler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(h.cfg.Session.Name)

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				utils.LogResponseError(c, err, h.logger)
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
			}
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(http.StatusInternalServerError, httpErrors.NewInternalServerError(err))
		}

		if err := h.sessionUsecase.DeleteSessionByID(c.Request().Context(), cookie.Value); err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		utils.DeleteSessionCookie(c, h.cfg.Session.Name)

		return c.NoContent(http.StatusOK)
	}
}

// Update godoc
// @Summary Update user
// @Description update existing user
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/{user_id} [put]
func (h *authHandler) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user := &models.User{}
		user.ID = uID

		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		updatedUser, err := h.authUsecase.UpdateUser(c.Request().Context(), user)
		if err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

// Delete
// @Summary Delete user account
// @Description Delete user account and remove cached data
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/{user_id} [delete]
func (h *authHandler) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err := h.authUsecase.DeleteUserByID(c.Request().Context(), uID); err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// GetUserByID godoc
// @Summary get user by ID
// @Description get user by ID
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param id path int true "user_id"
// @Success 200 {object} models.User
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/{user_id} [get]
func (h *authHandler) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user, err := h.authUsecase.GetUserByID(c.Request().Context(), uID)
		if err != nil {
			utils.LogResponseError(c, err, h.logger)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

// GetMe godoc
// @Summary Get the current logged in user
// @Description Get the current logged in user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/me [get]
func (h *authHandler) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok {
			utils.LogResponseError(c, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized), h.logger)
			return utils.ErrResponseWithLog(c, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		return c.JSON(http.StatusOK, user)
	}
}

// GetCSRFToken godoc
// @Summary Get CSRF token
// @Description Get CSRF token, required auth session cookie
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "Ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/token [get]
func (h *authHandler) GetCSRFToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		sid, ok := c.Get("sid").(string)
		if !ok {
			utils.LogResponseError(c, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized), h.logger)
			return utils.ErrResponseWithLog(c, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		token := csrf.MakeToken(sid, h.logger)
		c.Response().Header().Set(csrf.CSRFHeader, token)
		c.Response().Header().Set("Access-Control-Expose-Headers", csrf.CSRFHeader)

		return c.NoContent(http.StatusOK)
	}
}
