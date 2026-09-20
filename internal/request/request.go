package request

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/LuisCabantac/scholaflow-api/internal/apperrors"
)

const maxBodyBytes = 1_048_576

type contextKey string

const authUserContextKey contextKey = "authUser"

type AuthUser struct {
	ID    string
	Email string
	Name  string
}

func SetAuthUser(r *http.Request, user *AuthUser) *http.Request {
	ctx := context.WithValue(r.Context(), authUserContextKey, user)
	return r.WithContext(ctx)
}

func GetAuthUser(r *http.Request) (*AuthUser, bool) {
	user, ok := r.Context().Value(authUserContextKey).(*AuthUser)
	return user, ok && user != nil
}

type Validator interface {
	Validate() error
}

func DecodeJSON[T any](r *http.Request, dst *T) error {
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return &apperrors.AppError{
			Message:    "Malformed JSON request body.",
			Code:       "malformed_json",
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return &apperrors.AppError{
			Message:    "Request body must contain only a single JSON object.",
			Code:       "multiple_json_objects",
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	if v, ok := any(dst).(Validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}
