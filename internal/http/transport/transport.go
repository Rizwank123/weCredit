package transport

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/leebenson/conform"

	"github.com/weCredit/internal/domain"
)

const PageMax = 500

// DecodeAndValidateRequestBody decodes and validates the request body
func DecodeAndValidateRequestBody(ctx echo.Context, t interface{}) error {
	err := ctx.Bind(t)
	if err != nil {
		return err
	}
	err = conform.Strings(t)
	if err != nil {
		return err
	}
	err = ctx.Validate(t)
	if err != nil {
		return err
	}
	return nil
}

// SendResponse sends a response
func SendResponse(ctx echo.Context, status int, data interface{}) error {
	var finalResult domain.BaseResponse
	if data != nil {
		finalResult = domain.BaseResponse{
			Data: data,
		}
	}
	if status == http.StatusNoContent {
		return ctx.NoContent(status)
	}
	return ctx.JSON(status, finalResult)
}

// CustomValidator custom validator for echo
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate validates the request body
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
