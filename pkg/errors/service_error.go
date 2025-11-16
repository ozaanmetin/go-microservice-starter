package errors

import "fmt"

// Service error structure representing an error in API level operations, and its handled in middleware.

type ServiceError struct {
	// Internal fields
	StatusCode int   `json:"-"`
	Err        error `json:"-"`
	// Client-facing fields
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (s *ServiceError) Error() string {
	if s.Err != nil {
		return fmt.Sprintf("%s: %v", s.Message, s.Err)
	}
	return s.Message
}

// Unwrap returns the wrapped error for error wrapping support
func (s *ServiceError) Unwrap() error {
	return s.Err
}

func (s *ServiceError) AddDetail(key string, value interface{}) *ServiceError {
	if s.Details == nil {
		s.Details = make(map[string]interface{})
	}
	s.Details[key] = value
	return s
}

// Basic ServiceError implementations
func NewServiceError(statusCode int, code string, message string, err error) *ServiceError {
	return &ServiceError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Err:        err,
	}
}

func NewInternalServerError(err error) *ServiceError {
	return NewServiceError(500, "internal_server_error", "An internal server error occurred.", err)
}

func NewBadRequestError(message string, err error) *ServiceError {
	return NewServiceError(400, "bad_request", message, err)
}

func NewNotFoundError(message string, err error) *ServiceError {
	return NewServiceError(404, "not_found", message, err)
}

func NewUnauthorizedError(message string, err error) *ServiceError {
	return NewServiceError(401, "unauthorized", message, err)
}

func NewForbiddenError(message string, err error) *ServiceError {
	return NewServiceError(403, "forbidden", message, err)
}

func NewConflictError(message string, err error) *ServiceError {
	return NewServiceError(409, "conflict", message, err)
}

func NewTooManyRequestsError(message string, err error) *ServiceError {
	return NewServiceError(429, "too_many_requests", message, err)
}

func NewServiceUnavailableError(message string, err error) *ServiceError {
	return NewServiceError(503, "service_unavailable", message, err)
}
