package classrooms

import (
	"fmt"
	"net/http"

	"github.com/LuisCabantac/scholaflow-api/internal/apperrors"
	"github.com/LuisCabantac/scholaflow-api/internal/request"
	"github.com/LuisCabantac/scholaflow-api/internal/response"
)

type handler struct {
	service Service
}

func NewHandler(svc Service) *handler {
	return &handler{
		service: svc,
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var createClassroomReq CreateClassroomRequest
	err := request.DecodeJSON(r, &createClassroomReq)
	if err != nil {
		response.Error(w, fmt.Errorf("failed to parse create classroom payload: %v: %w", err, apperrors.ErrMissingBody))
		return
	}

	user, ok := request.GetAuthUser(r)
	if !ok {
		response.Error(w, apperrors.ErrUnauthorizedAccess)
		return
	}

	classroom, err := h.service.Create(r.Context(), createClassroomReq, user.ID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "Classroom created successfully.", classroom)
}

func (h *handler) Enroll(w http.ResponseWriter, r *http.Request) {
	var enrollClassroomReq EnrollClassroomRequest
	err := request.DecodeJSON(r, &enrollClassroomReq)
	if err != nil {
		response.Error(w, fmt.Errorf("failed to parse classroom enroll payload: %v: %w", err, apperrors.ErrMissingBody))
		return
	}

	user, ok := request.GetAuthUser(r)
	if !ok {
		response.Error(w, apperrors.ErrUnauthorizedAccess)
		return
	}

	classroomEnrollment, err := h.service.Enroll(r.Context(), enrollClassroomReq, user.ID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "Enrolled to classroom successfully.", classroomEnrollment)
}
