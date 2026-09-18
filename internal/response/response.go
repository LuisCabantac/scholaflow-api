package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/LuisCabantac/scholaflow-api/internal/apperrors"
)

func JSON[T comparable](w http.ResponseWriter, status int, data T) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		slog.Error("failed to encode JSON", "error", err.Error())
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(buf.Bytes())
	if err != nil {
		slog.Error("failed to write response", "error", err.Error())
		return err
	}

	return nil
}

func Text(w http.ResponseWriter, status int, body string) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, err := w.Write([]byte(body))
	return err
}

func Error(w http.ResponseWriter, err error) error {
	slog.Error("request failed", "error", err.Error())

	if appErr, ok := errors.AsType[*apperrors.AppError](err); ok {
		return JSON(w, appErr.StatusCode, appErr)
	}

	return JSON(w, http.StatusInternalServerError, apperrors.ErrInternalServerError)
}
