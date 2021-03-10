package handler

import (
	"bitbucket.org/latonaio/authenticator/internal/crypto"
	"bitbucket.org/latonaio/authenticator/internal/models"
	custmres "bitbucket.org/latonaio/authenticator/pkg/response"
	"github.com/labstack/echo/v4"
)

type RegisterUserParam struct {
	LoginID  string `json:"login_id" form:"login_id"`
	Password string `json:"password" form:"password"`
}

func RegisterUser(c echo.Context) error {
	param := &RegisterUserParam{}
	err := c.Bind(param)
	if err != nil {
		return c.JSON(custmres.BadRequestRes.Code, custmres.BadRequestRes)
	}

	// validate input fields.
	unverifiedUser := &models.User{
		LoginID:  param.LoginID,
		Password: param.Password,
	}
	if err := unverifiedUser.Validate(); err != nil {
		c.Logger().Printf("Failed to validate input parameter: %v", err)
		return c.JSON(custmres.BadRequestRes.Code, custmres.BadRequestRes)
	}

	// check registration status of login id.
	user := models.NewUser()
	result, err := user.GetByLoginID(param.LoginID)
	if result != nil && err == nil {
		c.Logger().Printf("Login id is already used")
		return c.JSON(custmres.Conflict.Code, custmres.Conflict.Message)
	}

	// encrypt password
	encryptedPassword, err := crypto.Encrypt(param.Password)
	if err != nil {
		c.Logger().Printf("Failed to encrypt password: %v", err)
		return c.JSON(custmres.InternalErrRes.Code, custmres.InternalErrRes)
	}
	user.SetUser(&models.User{
		LoginID:  param.LoginID,
		Password: encryptedPassword,
	})

	err = user.Register()
	if err != nil {
		return err
	}
	return nil
}
