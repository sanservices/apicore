package apierror

const (
	// CodeNotFoundErr is used when a resource is not found.
	CodeNotFoundErr = "account/not-found"

	// CodePropertyValidationErr when a validation rule for a property fails.
	CodePropertyValidationErr = "validation/property"

	// CodeFormatNotValidErr when a payload is in the incorrect format.
	CodeFormatNotValidErr = "validation/incorrect-format"

	// CodeDecodingErr when decoding a payload fails.
	CodeDecodingErr = "validation/decoding"

	// CodeInternalErr for internal server errors.
	CodeInternalErr = "service/internal"

	// CodeAccountNotAvailableErr when trying to use an account not available.
	CodeAccountNotAvailableErr = "account/not-available"
)

func New(text string, code string) error {
	return &ApiError{s: text, c: code}
}

// ApiError is the struct for the custom error.
type ApiError struct {
	s string
	c string
}

func (e *ApiError) Error() string {
	return e.s
}

func (e *ApiError) Code() string {
	return e.c
}
