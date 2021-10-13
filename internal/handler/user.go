package handler

import (
	"errors"
	"net/http"
	"time"

	"bitbucket.org/latonaio/authenticator/internal/crypto"
	"bitbucket.org/latonaio/authenticator/internal/models"
	custmres "bitbucket.org/latonaio/authenticator/pkg/response"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserParam struct {
	LoginID   string `json:"login_id" form:"login_id"`
	Password  string `json:"password" form:"password"`
	Qos       string `json:"qos" form:"qos"`
	IsEncrypt bool   `json:"is_encrypt" form:"is_encrypt"`
}

func RegisterUser(c echo.Context) error {
	param := &UserParam{}
	param.IsEncrypt = true // default value
	err := c.Bind(param)
	if err != nil {
		return c.JSON(custmres.BadRequestRes.Code, custmres.BadRequestRes)
	}

	// validate input fields.
	unverifiedUser := &models.User{
		LoginID:  param.LoginID,
		Password: param.Password,
		Qos:      models.ToQos(param.Qos),
	}
	if unverifiedUser.NeedsValidation() {
		if err := unverifiedUser.Validate(); err != nil {
			c.Logger().Printf("Failed to validate input parameter: %v", err)
			return c.JSON(custmres.BadRequestRes.Code, custmres.BadRequestRes)
		}
	}

	// check registration status of login id.
	user := models.NewUser()
	result, err := user.GetByLoginID(param.LoginID)
	if result != nil && err == nil {
		c.Logger().Printf("Login id is already used")
		return c.JSON(custmres.Conflict.Code, custmres.Conflict.Message)
	}

	userImp := &models.User{
		LoginID:     param.LoginID,
		Password:    param.Password,
		Qos:         unverifiedUser.Qos,
		IsEncrypt:   &param.IsEncrypt,
		LastLoginAt: nil,
	}

	// encrypt password
	if param.IsEncrypt {
		encryptedPassword, err := crypto.Encrypt(param.Password)
		if err != nil {
			c.Logger().Printf("Failed to encrypt password: %v", err)
			return c.JSON(custmres.InternalErrRes.Code, custmres.InternalErrRes)
		}
		userImp.Password = encryptedPassword
	}
	user.SetUser(userImp)

	err = user.Register()
	if err != nil {
		return err
	}
	return nil
}

type UpdateUserParam struct {
	UserParam
	OldPassword string `json:"old_password" form:"old_password"`
}

func (p *UpdateUserParam) PasswordExists() bool {
	return p.Password != ""
}

func (p *UpdateUserParam) QosExists() bool {
	return p.Qos != ""
}

func UpdateUser(c echo.Context) error {
	param := &UpdateUserParam{}
	param.IsEncrypt = true // default value
	err := c.Bind(param)
	if err != nil {
		return c.JSON(custmres.BadRequestRes.Code, custmres.BadRequestRes)
	}

	// check existence of user.
	result, err := models.NewUser().GetByLoginID(c.Param("login_id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Logger().Printf("Login id is not found: %s", param.LoginID)
			return c.JSON(custmres.NotFoundErrRes.Code, custmres.NotFoundErrRes)
		} else {
			c.Logger().Printf("Failed to db access: %v", err)
			return c.JSON(custmres.InternalErrRes.Code, custmres.InternalErrRes)
		}
	}
	if result != nil && result.IsDeleted() {
		c.Logger().Print("User is already deleted.")
		return c.JSON(custmres.Conflict.Code, custmres.Conflict.Message)
	}

	// authenticate old password.
	if !*result.IsEncrypt {
		if result.Password != param.OldPassword {
			c.Logger().Print("Failed to login due to incorrect password")
			return c.JSON(custmres.UnauthorizedRes.Code, custmres.UnauthorizedRes)
		}
	} else {
		if err := crypto.CompareHashAndPassword(result.Password, param.OldPassword); err != nil {
			c.Logger().Printf("Failed to login: %v", err)
			return c.JSON(custmres.UnauthorizedRes.Code, custmres.UnauthorizedRes)
		}
	}

	userImp := &models.User{
		LoginID:   param.LoginID,
		Password:  param.Password,
		Qos:       models.ToQos(param.Qos),
		IsEncrypt: &param.IsEncrypt,
	}
	if !param.QosExists() {
		userImp.Qos = result.Qos
	}

	// validate input params.
	if userImp.NeedsValidation() {
		if err := userImp.Validate(); err != nil {
			c.Logger().Printf("Failed to validate input parameter: %v", err)
			return c.JSON(custmres.BadRequestRes.Code, custmres.BadRequestRes)
		}
	}

	// encrypt password.
	if param.PasswordExists() && param.IsEncrypt {
		encryptedPassword, err := crypto.Encrypt(param.Password)
		if err != nil {
			c.Logger().Printf("Failed to encrypt password: %v", err)
			return c.JSON(custmres.InternalErrRes.Code, custmres.InternalErrRes)
		}
		userImp.Password = encryptedPassword
	}

	result.SetUser(userImp)
	if err = result.Update(); err != nil {
		return err
	}
	return nil
}

func GetUser(c echo.Context) error {
	loginID := c.Param("login_id")
	result, err := models.NewUser().GetByLoginID(loginID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Logger().Printf("Login id is not found: %s", loginID)
			return c.JSON(custmres.NotFoundErrRes.Code, custmres.NotFoundErrRes)
		} else {
			c.Logger().Printf("Failed to db access: %v", err)
			return c.JSON(custmres.InternalErrRes.Code, custmres.InternalErrRes)
		}
	}
	if result != nil && result.IsDeleted() {
		c.Logger().Print("User is already deleted.")
		return c.JSON(custmres.Conflict.Code, custmres.Conflict.Message)
	}
	return c.JSON(http.StatusOK, custmres.UserResponseFormat{LoginID: result.LoginID})
}

type DeleteUserParam struct {
	Password string `json:"password" form:"password"`
}

func DeleteUser(c echo.Context) error {

	// ユーザの理論削除を行います。
	// 具体的な処理としては、ユーザモデルの deleted_at フラグに現在時刻を登録します。

	param := &DeleteUserParam{}
	err := c.Bind(param)
	if err != nil {
		return c.JSON(custmres.BadRequestRes.Code, custmres.BadRequestRes)
	}
	loginID := c.Param("login_id")

	result, err := models.NewUser().GetByLoginID(loginID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Logger().Printf("Login id is not found: %s", loginID)
			return c.JSON(custmres.NotFoundErrRes.Code, custmres.NotFoundErrRes)
		} else {
			c.Logger().Printf("Failed to db access: %v", err)
			return c.JSON(custmres.InternalErrRes.Code, custmres.InternalErrRes)
		}
	}
	if result != nil && result.IsDeleted() {
		c.Logger().Print("User is already deleted.")
		return c.JSON(custmres.Conflict.Code, custmres.Conflict.Message)
	}

	// authenticate password.
	if !*result.IsEncrypt {
		if result.Password != param.Password {
			c.Logger().Print("Failed to login due to incorrect password")
			return c.JSON(custmres.UnauthorizedRes.Code, custmres.UnauthorizedRes)
		}
	} else {
		if err := crypto.CompareHashAndPassword(result.Password, param.Password); err != nil {
			c.Logger().Printf("Failed to login: %v", err)
			return c.JSON(custmres.UnauthorizedRes.Code, custmres.UnauthorizedRes)
		}
	}
	now := time.Now()
	userImp := &models.User{ID: result.ID, DeletedAt: &now}
	if err := userImp.Update(); err != nil {
		return err
	}
	return nil
}
