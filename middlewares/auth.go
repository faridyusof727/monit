package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	AuthClient *auth.Client
}

func (a *Auth) Check(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		token = strings.Replace(token, "Bearer ", "", 1)

		user, err := a.AuthClient.VerifyIDToken(context.Background(), token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		c.Request().Header.Set("UID", user.UID)
		c.Request().Header.Set("email", fmt.Sprint(user.Claims["email"]))

		return next(c)
	}
}
