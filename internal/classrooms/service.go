package classrooms

import (
	"context"
	"errors"
	"math/rand/v2"

	"github.com/LuisCabantac/scholaflow-api/internal/adapters/postgresql"
	postgres "github.com/LuisCabantac/scholaflow-api/internal/adapters/postgresql/sqlc"
	"github.com/LuisCabantac/scholaflow-api/internal/apperrors"
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

func (s *svc) CreateClassroom(ctx context.Context, req CreateClassroomRequest, userID string) (*postgres.Classroom, error) {
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

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return &classroom, nil
}
