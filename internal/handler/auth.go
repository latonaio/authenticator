package handler

import (
	"net/http"
	"os"
	"time"

	"bitbucket.org/latonaio/authenticator/internal/crypto"
	"bitbucket.org/latonaio/authenticator/internal/models"
	custmerr "bitbucket.org/latonaio/authenticator/pkg/error"
	custmres "bitbucket.org/latonaio/authenticator/pkg/response"
	"github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo/v4"
)

type UserLoginParam struct {
	//　TODO: アカウント作成時は、"user_id" と "password" なので、どちらかに統一したい。
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
	if err := crypto.CompareHashAndPassword(result.Password, param.Password); err != nil {
		c.Logger().Printf("Failed to login: %v", err)
		return c.String(http.StatusUnauthorized, "The login ID or password you entered was incorrect")
	}

	// generate JWT
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.User().ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // TODO: configs.yml とか環境変数とかに記述する
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("PRIVATE_KEY")))
	if err != nil {
		c.Logger().Printf("Failed to parse private key: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to generate access token.")
	}
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		c.Logger().Printf("Failed to generate JWT: %v", err)
		return c.String(http.StatusInternalServerError, "Failed to generate access token.")
	}

	return c.String(http.StatusOK, signedToken)
}
