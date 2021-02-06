package route

import (
	"bitbucket.org/latonaio/authenticator/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterRouting(e *echo.Echo) {
	// ユーザー作成
	e.POST("/users", handler.RegisterUser)
	// ユーザー認証
	e.POST("/login", handler.EnsureUser)
}
