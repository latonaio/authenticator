package handler

import (
	"net/http"
	"os"
	"time"

	"bitbucket.org/latonaio/authenticator/configs"
	"bitbucket.org/latonaio/authenticator/internal/crypto"
	"bitbucket.org/latonaio/authenticator/internal/models"
	custmres "bitbucket.org/latonaio/authenticator/pkg/response"
	"github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo/v4"
)

type UserLoginParam struct {
	LoginID  string `json:"login_id" form:"login_id"`
	Password string `json:"password" form:"password"`
}

var jwtExp int64

func init() {
	cfgs, err := configs.New()
	if err != nil {
		panic(err)
	}
	jwtExp = cfgs.Get().Jwt.Exp
}

func EnsureUser(c echo.Context) error {
	param := &UserLoginParam{}
	err := c.Bind(param)
	if err != nil {
		return c.JSON(custmres.BadRequestRes.Code, custmres.BadRequestRes)
	}
	user := models.NewUser()
	result, err := user.GetByLoginID(param.LoginID)
	if err != nil {
		return c.JSON(custmres.NotFoundErrRes.Code, custmres.NotFoundErrRes)
	}
	if err := crypto.CompareHashAndPassword(result.Password, param.Password); err != nil {
		c.Logger().Printf("Failed to login: %v", err)
		return c.JSON(custmres.UnauthorizedRes.Code, custmres.UnauthorizedRes)
	}

	// generate JWT
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.User().ID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(jwtExp)).Unix()
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("PRIVATE_KEY")))
	if err != nil {
		c.Logger().Printf("Failed to parse private key: %v", err)
		return c.JSON(custmres.InternalErrRes.Code, custmres.InternalErrRes)
	}
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		c.Logger().Printf("Failed to generate JWT: %v", err)
		return c.JSON(custmres.InternalErrRes.Code, custmres.InternalErrRes)
	}

	if err := result.Login(); err != nil {
		c.Logger().Printf("Failed to record last_login_at: %v", err)
		return c.JSON(custmres.InternalErrRes.Code, custmres.InternalErrRes)
	}
	return c.JSON(http.StatusOK, custmres.JWTResponseFormat{Jwt: signedToken})
}
