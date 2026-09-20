package classrooms

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/LuisCabantac/scholaflow-api/internal/adapters/postgresql"
	postgres "github.com/LuisCabantac/scholaflow-api/internal/adapters/postgresql/sqlc"
	"github.com/LuisCabantac/scholaflow-api/internal/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func generateClassCode() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

type svc struct {
	queries *postgres.Queries
}

func NewService(q *postgres.Queries) *svc {
	return &svc{
		queries: q,
	}
}

const maxRetries = 3

func (s *svc) Create(ctx context.Context, req CreateClassroomRequest, userID string) (*postgres.Classroom, error) {
	for range maxRetries {
		classroom, err := s.queries.CreateClassroom(ctx, postgres.CreateClassroomParams{
			Name:              req.Name,
			Subject:           postgresql.TextFromPtr(req.Subject),
			Section:           req.Section,
			Description:       postgresql.TextFromPtr(req.Description),
			Room:              postgresql.TextFromPtr(req.Room),
			Code:              generateClassCode(),
			CardBackground:    req.CardBackground,
			IllustrationIndex: int32(rand.IntN(5)),
			TeacherID:         userID,
		})

		if err == nil {
			return &classroom, nil
		}

		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23503" {
				return nil, apperrors.ErrUserNotFound
			}

			if pgErr.Code == "23505" && pgErr.ConstraintName == "classroom_code_key" {
				continue
			}
		}

		return nil, err
	}

	return nil, ErrCodeGenerationFailed
}

func (s *svc) Enroll(ctx context.Context, req EnrollClassroomRequest, userID string) (*postgres.ClassroomEnrollment, error) {
	classroom, err := s.queries.GetClassroomByCode(ctx, req.Code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClassNotFound
		}

		return nil, fmt.Errorf("failed to fetch classroom: %w", err)
	}

	if classroom.TeacherID == userID {
		return nil, ErrSelfEnrollment
	}

	enrollment, err := s.queries.CreateClassroomEnrollment(ctx, postgres.CreateClassroomEnrollmentParams{
		ClassroomID: classroom.ID,
		UserID:      userID,
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "uq_classroom_enrollment" {
				return nil, ErrAlreadyEnrolled
			}
		}
		return nil, fmt.Errorf("faild to enroll in classroom: %w", err)
	}

	return &enrollment, err

}
