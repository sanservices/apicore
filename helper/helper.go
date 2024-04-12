package helper

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sanservices/apicore/errs"
)

// httpResponse standard json response
type successResponse struct {
	Data interface{} `json:"data"`
}

// JSON response with error
type errorResponse struct {
	Errors []errs.ServiceError `json:"errors"`
}

// RespondError sends a json error response
func RespondError(c echo.Context, httpStatusCode int, respErr error) error {
	var response errorResponse
	if sErrs, ok := respErr.(errs.ServiceError); ok {
		response.Errors = append(response.Errors, sErrs)
	} else {
		switch httpStatusCode {
		case http.StatusBadRequest:
			response.Errors = append(response.Errors, errs.NewInputMalformedErr(respErr))
		default:
			response.Errors = append(response.Errors, errs.NewInternalErr(respErr))
		}
	}

	return c.JSON(httpStatusCode, response)
}

// RespondOk sends a json success response
func RespondOk(c echo.Context, data interface{}) error {
	response := &successResponse{Data: data}

	return c.JSON(http.StatusOK, response)
}

// DecodeParams decodes the query params into destination interface
func DecodeQueryParams(c echo.Context, dst interface{}) error {
	binder := &echo.DefaultBinder{}
	return binder.BindQueryParams(c, dst)
}

// DecodePathParams decodes the path params into destination interface
func DecodePathParams(c echo.Context, dst interface{}) error {
	binder := &echo.DefaultBinder{}
	return binder.BindPathParams(c, dst)
}

// DecodeBody decodes the request body into destination interface
func DecodeBody(c echo.Context, dst interface{}) error {
	data := c.Request().Body
	decoder := json.NewDecoder(data)
	return decoder.Decode(dst)
}
