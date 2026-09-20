package classrooms

import (
	"net/http"

	"github.com/LuisCabantac/scholaflow-api/internal/apperrors"
)

var (
	ErrCodeGenerationFailed = &apperrors.AppError{
		Message:    "Failed to generate a unique classroom code. Please try again.",
		Code:       "code_generation_failed",
		StatusCode: http.StatusInternalServerError}
	ErrClassNotFound = &apperrors.AppError{
		Message:    "Class not found.",
		Code:       "class_not_found",
		StatusCode: http.StatusNotFound}
	ErrSelfEnrollment = &apperrors.AppError{
		Message:    "You cannot enroll in your own classroom.",
		Code:       "cannot_enroll_own_class",
		StatusCode: http.StatusBadRequest}
	ErrAlreadyEnrolled = &apperrors.AppError{
		Message:    "You are already enrolled in this classroom.",
		Code:       "already_enrolled",
		StatusCode: http.StatusConflict}
	ErrNotEnrolled = &apperrors.AppError{
		Message:    "You are not enrolled in this classroom.",
		Code:       "not_enrolled",
		StatusCode: http.StatusConflict}
)
