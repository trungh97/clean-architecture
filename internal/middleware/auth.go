package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"instagram-clone.com/m/pkg/httpErrors"
	"instagram-clone.com/m/pkg/utils"
)

// Auth session middleware using redis
func (mw *MiddlewareManager) AuthSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(mw.cfg.Session.Name)

		if err != nil {
			mw.logger.Errorf("AuthSessionMiddleware RequestID: %s, Error: %v", utils.GetRequestID(c), err.Error())
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
			}
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		sid := cookie.Value
		session, err := mw.sessionUsecase.GetSessionByID(c.Request().Context(), sid)
		if err != nil {
			mw.logger.Errorf("GetSessionByID RequestID: %s, CookieValue: %s, Error: %v", utils.GetRequestID(c), sid, err.Error())
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		user, err := mw.authUsecase.GetUserByID(c.Request().Context(), session.UserID)
		if err != nil {
			mw.logger.Errorf("GetUserByID RequestID: %s, SessionID: %s, Error: %v", utils.GetRequestID(c), sid, err.Error())
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		c.Set("sid", sid)
		c.Set("user", user)
		c.Set("uid", session.SessionID)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, user)
		c.SetRequest(c.Request().WithContext(ctx))

		mw.logger.Info("SessionMiddleware, RequestID: %s, IP: %s, UserID: %s, CookieSessionID: %s", utils.GetRequestID(c), utils.GetIPAddress(c), user.ID.String(), cookie.Value)

		return next(c)
	}
}
