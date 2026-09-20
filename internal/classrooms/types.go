package classrooms

import (
	"context"
	"net/http"
	"strings"

	postgres "github.com/LuisCabantac/scholaflow-api/internal/adapters/postgresql/sqlc"
	"github.com/LuisCabantac/scholaflow-api/internal/apperrors"
)

type Service interface {
	CreateClassroom(ctx context.Context, req CreateClassroomRequest, userID string) (*postgres.Classroom, error)
}

type CreateClassroomRequest struct {
	Name           string  `json:"name"`
	Subject        *string `json:"subject,omitempty"`
	Section        string  `json:"section"`
	Description    *string `json:"description,omitempty"`
	Room           *string `json:"room,omitempty"`
	CardBackground string  `json:"cardBackground"`
}

func (r *CreateClassroomRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return &apperrors.AppError{
			Message:    "Classroom name is required",
			Code:       "missing_name",
			StatusCode: http.StatusBadRequest,
		}
	}

	if strings.TrimSpace(r.Section) == "" {
		return &apperrors.AppError{
			Message:    "Classroom section is required",
			Code:       "missing_section",
			StatusCode: http.StatusBadRequest,
		}
	}

	if strings.TrimSpace(r.CardBackground) == "" {
		return &apperrors.AppError{
			Message:    "Classroom card background is required",
			Code:       "missing_card_background",
			StatusCode: http.StatusBadRequest,
		}
	}

	if r.Subject != nil && strings.TrimSpace(*r.Subject) == "" {
		r.Subject = nil
	}

	if r.Description != nil && strings.TrimSpace(*r.Description) == "" {
		r.Description = nil
	}

	if r.Room != nil && strings.TrimSpace(*r.Room) == "" {
		r.Room = nil
	}

	return nil
}
