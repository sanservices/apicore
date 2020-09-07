package apicore

type ResponsePropertyError struct {
	Property    string   `json:"property"`
	Constraints []string `json:"constraints"`
}

type ResponseError struct {
	Message    string                  `json:"message"`
	Code       string                  `json:"code"`
	Properties []ResponsePropertyError `json:"properties"`
}

type ResponseErr interface {
	Error() string
}

func (re ResponseError) Error() string {
	return re.Message
}
