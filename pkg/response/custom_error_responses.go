package custmres

import (
	"net/http"

	customerror "bitbucket.org/latonaio/authenticator/pkg/error"
	"github.com/labstack/echo/v4"
)

type ResponseFormat struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	res := generateErrorResponse(err)
	_ = c.JSON(res.Code, &res)
	c.Logger().Error(err)
}

func generateErrorResponse(err error) ResponseFormat {
	custmErr, ok := err.(customerror.CustomErrMessage)
	if !ok {
		return InternalErrRes
	}
	switch custmErr {
	case customerror.ErrBadRequest:
		return BadRequestRes
	case customerror.ErrNotFound:
		return NotFoundErrRes
	}
	return InternalErrRes
}

var (
	BadRequestRes = ResponseFormat{
		Code:    http.StatusBadRequest,
		Message: customerror.ErrBadRequest.Error(),
	}

	InternalErrRes = ResponseFormat{
		Code:    http.StatusInternalServerError,
		Message: customerror.ErrInternal.Error(),
	}

	NotFoundErrRes = ResponseFormat{
		Code:    http.StatusNotFound,
		Message: customerror.ErrNotFound.Error(),
	}

	UnauthorizedRes = ResponseFormat{
		Code:    http.StatusUnauthorized,
		Message: customerror.ErrUnauthorized.Error(),
	}

	Conflict = ResponseFormat{
		Code:    http.StatusConflict,
		Message: customerror.ErrConflict.Error(),
	}
)
