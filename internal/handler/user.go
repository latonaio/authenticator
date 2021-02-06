package handler

import (
	"bitbucket.org/latonaio/authenticator/internal/models"
	custmres "bitbucket.org/latonaio/authenticator/pkg/response"
	"github.com/labstack/echo/v4"
)

type RegisterUserParam struct {
	UserID   string `json:"user_id" form:"user_id"`
	Password string `json:"user_password" form:"user_password"`
}

func RegisterUser(c echo.Context) error {
	param := &RegisterUserParam{}
	err := c.Bind(param)
	if err != nil {
		return c.JSON(custmres.BadRequestRes.Code, custmres.BadRequestRes)
	}
	user := models.NewUser()

	// TODO password hash化
	user.SetUser(&models.User{
		LoginID:  param.UserID,
		Password: param.Password,
	})
	err = user.Register()
	if err != nil {
		return err
	}
	return nil
}