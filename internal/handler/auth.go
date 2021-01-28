package handler

import (
	"net/http"

	"bitbucket.org/latonaio/authenticator/internal/models"
	custmerr "bitbucket.org/latonaio/authenticator/pkg/error"
	custmres "bitbucket.org/latonaio/authenticator/pkg/response"
	"github.com/labstack/echo/v4"
)

type UserLoginParam struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"user_password" form:"user_password"`
}

func EnsureUser(c echo.Context) error {
	param := &UserLoginParam{}
	err := c.Bind(param)
	if err != nil {
		return c.JSON(custmres.BadRequestRes.Code, custmres.BadRequestRes)
	}
	user := models.NewUser()
	result, err := user.GetByLoginID(param.UserName)
	if err != nil {
		return custmerr.ErrNotFound
	}

	// TODO verify password hash
	if result.Password != param.Password {
		return c.String(http.StatusUnauthorized, "The login ID or password you entered was incorrect")
	}

	return c.String(http.StatusOK, "success")
}
