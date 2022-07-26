package alma

import "fmt"

// UnauthorizedError occurs when API key is wrong or IP is banned.
type UnauthorizedError struct {
	errorMessage string
}

func (e *UnauthorizedError) Error() string {
	return e.errorMessage
}

// InvalidRequestError occurs when the request is ill-formed : missing
// parameter, unknown identifier and so on.
type InvalidRequestError struct {
	errorMessage string
}

func (e *InvalidRequestError) Error() string {
	return e.errorMessage
}

// ServerError occurs when Alma responds with a 5XX error.
type ServerError struct {
	errorMessage string
}

func (e *ServerError) Error() string {
	return e.errorMessage
}

// NotFoundError occurs when the requested resource does not exist.
type NotFoundError struct {
	id           string
	errorMessage string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("id %s: %v", e.id, e.errorMessage)
}

// ThresholdError occurs when the number of concurrent requests hits the Alma
// limits : 200,000 requests/day and 25 requests/second.
type ThresholdError struct {
	errorMessage string
}

func (e *ThresholdError) Error() string {
	return e.errorMessage
}

// FetchError is used for any other error.
type FetchError struct {
	errorMessage string
}

func (e *FetchError) Error() string {
	return e.errorMessage
}
