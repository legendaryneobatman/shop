package webtool

type APIError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Source  error  `json:"source,omitempty"`
}

func NewAPIError(message string, code int, source error) *APIError {
	return &APIError{
		Code: code,
		Message: message,
		Source: source,
	}
}

func MakeError(message string, code int) func(err error) *APIError {
	return func (err error) *APIError {
		return NewAPIError(message, code, err)
	}
}