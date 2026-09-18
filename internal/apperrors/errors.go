package apperrors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Message    string `json:"message"`
	Code       string `json:"code,omitempty"`
	StatusCode int    `json:"statusCode"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrInternalServerError = &AppError{
		Message:    "Internal server error",
		Code:       "internal_error",
		StatusCode: http.StatusInternalServerError}
	ErrInvalidToken = &AppError{
		Message:    "Invalid or expired token.",
		Code:       "invalid_token",
		StatusCode: http.StatusUnauthorized}
	ErrMissingIDParam = &AppError{
		Message:    "ID parameter is required.",
		Code:       "missing_id_param",
		StatusCode: http.StatusBadRequest}
	ErrMissingBody = &AppError{
		Message:    "Request body is required.",
		Code:       "missing_request_body",
		StatusCode: http.StatusBadRequest}
	ErrUserNotFound = &AppError{
		Message:    "User not found.",
		Code:       "user_not_found",
		StatusCode: http.StatusNotFound}
	ErrUnauthorizedAccess = &AppError{
		Message:    "Forbidden.",
		Code:       "access_denied",
		StatusCode: http.StatusForbidden}
	ErrBadRequest = &AppError{
		Message:    "Bad request.",
		Code:       "bad_request",
		StatusCode: http.StatusBadRequest}
)

func BadRequest(message string, code ...string) *AppError {
	c := "bad_request"
	if len(code) > 0 {
		c = code[0]
	}
	return &AppError{
		Message:    message,
		Code:       c,
		StatusCode: http.StatusBadRequest,
	}
}

func NotFound(message string, code ...string) *AppError {
	c := "not_found"
	if len(code) > 0 {
		c = code[0]
	}
	return &AppError{
		Message:    message,
		Code:       c,
		StatusCode: http.StatusNotFound,
	}
}

func Internal(err error) *AppError {
	return &AppError{
		Message:    "Internal server error",
		Code:       "internal_error",
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
