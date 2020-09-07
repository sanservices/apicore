package apicore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator"
	"github.com/gorilla/schema"
	"github.com/san-services/apicore/apierror"
	"github.com/san-services/apilogger"
)

type ConsumerErrors struct {
	Errors []ResponseError `json:"errors"`
}

// httpRequest creates request struct
type httpRequest struct {
	Response *http.Response
	URL      string
	Data     []byte
	Method   string
	request  *http.Request
	Headers  map[string]string
}

// Send sends http request
func (r *httpRequest) Send() (err error) {
	httpClient := http.Client{}

	r.request, err = http.NewRequest(r.Method, r.URL, bytes.NewBuffer(r.Data))
	if err != nil {
		return err
	}

	// Add headers to the request
	for index, header := range r.Headers {
		r.request.Header.Add(index, header)
	}

	// make requests and store response
	r.Response, err = httpClient.Do(r.request)
	return
}

// httpResponse standard json response
type httpResponse struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// camelCaseSplit source code belongs to from https://github.com/fatih/camelcase/blob/master/camelcase.go
func camelCaseSplit(src string) (entries []string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return []string{src}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	class := 0

	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return
}

// DecodeBody decodes a json request body into a given struct
func DecodeBody(ctx context.Context, data io.ReadCloser, model interface{}) (err error) {
	err = json.NewDecoder(data).Decode(model)
	if err != nil {
		err = apierror.New(err.Error(), apierror.CodePropertyValidationErr)
	}

	return err
}

func dateParamConverter(value string) reflect.Value {
	v, err := time.Parse("2006-01-02", value)

	if err != nil {
		return reflect.Value{}
	}

	return reflect.ValueOf(v)
}

func DecodeParams(ctx context.Context, dst interface{}, src map[string][]string) (errors schema.MultiError) {
	d := schema.NewDecoder()
	d.SetAliasTag("json")
	d.RegisterConverter(time.Time{}, dateParamConverter)
	logger := apilogger.New(ctx, "")

	if err := d.Decode(dst, src); err != nil {
		errors, ok := err.(schema.MultiError)
		if !ok {
			logger.Errorf(apilogger.LogCatInputValidation, "Could not get errors %v", err)
		}
		return errors
	}

	return nil
}

// Respond sends standard json response
func Respond(ctx context.Context,
	w http.ResponseWriter, data interface{}, httpStatusCode int, err error) {

	if err != nil {
		RespondError(ctx, w, httpStatusCode, err)
		return
	}
	RespondSuccess(ctx, w, data)
}

// RespondSuccess sends a json success response
func RespondSuccess(
	ctx context.Context, w http.ResponseWriter, data interface{}) {

	response := &httpResponse{Data: data}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RespondError sends a json error response
func RespondError(
	ctx context.Context, w http.ResponseWriter, httpStatusCode int, err error) {

	// Get error code.
	var errCode string = apierror.CodeInternalErr
	apiErr, ok := err.(*apierror.ApiError)

	if ok {
		errCode = apiErr.Code()
	}

	// This condition is here to avoid
	// having to check for each error code
	// before passing it.
	if errCode == apierror.CodeNotFoundErr {
		httpStatusCode = http.StatusNotFound
	}

	// Base headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	// Define basic error to respond struct.
	errToRespond := &ResponseError{
		Message:    err.Error(),
		Code:       errCode,
		Properties: []ResponsePropertyError{},
	}

	// Check if the error was already set.
	respErr, ok := err.(*ResponseError)
	if ok {
		errToRespond = respErr
	}

	// Check if the error comes from validating payload.
	errors, ok := err.(validator.ValidationErrors)

	if ok {
		errToRespond.Message = "Invalid parameters in request"
		errToRespond.Code = apierror.CodePropertyValidationErr

		for _, validationErr := range errors {
			var constraint string
			if validationErr.Param() != "" {
				constraint = fmt.Sprintf("validation [%s %s] failed with value: %v", validationErr.Tag(), validationErr.Param(), validationErr.Value())
			} else {
				constraint = fmt.Sprintf("validation [%s] failed with value: %v", validationErr.Tag(), validationErr.Value())
			}

			// This code is use because validator package
			// doesn't offer a method to get the json property name.
			propertySplitted := camelCaseSplit(validationErr.Field())
			propertyJoined := strings.Join(propertySplitted, "_")
			property := strings.ToLower(propertyJoined)

			errToRespond.Properties = append(errToRespond.Properties, ResponsePropertyError{
				Property:    property,
				Constraints: []string{constraint},
			})
		}
	}

	errs := ConsumerErrors{[]ResponseError{*errToRespond}}
	json.NewEncoder(w).Encode(errs)
}
