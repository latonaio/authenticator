package route

import (
	"bitbucket.org/latonaio/authenticator/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterRouting(e *echo.Echo) {
	e.POST("/users", handler.RegisterUser)
	e.POST("/login", handler.EnsureUser)
}
