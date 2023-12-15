package repository

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewRepository(db *sqlx.DB, tracer trace.Tracer) *Repository {
	return &Repository{db: db, tracer: tracer}
}

func (repo *Repository) Create(ctx context.Context, userID int, input domain.Quiz) (int, error) {
	ctx, span := repo.tracer.Start(ctx, "admin.quizRepo.Create")
	defer span.End()

	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO quizzes (title, description, user_id) VALUES ($1, $2, $3) RETURNING id",
		input.Title, input.Description, userID).Scan(&id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetQuizzes(ctx context.Context) ([]models.Quiz, error) {
	ctx, span := repo.tracer.Start(ctx, "admin.quizRepo.GetQuizzes")
	defer span.End()

	quizzes := make([]models.Quiz, 0)

	err := repo.db.SelectContext(ctx, &quizzes, "SELECT id, title, description, image, user_id, created_at FROM quizzes")

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return quizzes, nil
}

func (repo *Repository) Update(ctx context.Context, id int, input domain.Quiz) error {
	ctx, span := repo.tracer.Start(ctx, "admin.quizRepo.Update")
	defer span.End()

	err := repo.db.QueryRowxContext(ctx, "UPDATE quizzes SET title = $1, description = $2 WHERE id = $3",
		input.Title, input.Description, id).Err()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, id int) error {
	ctx, span := repo.tracer.Start(ctx, "admin.quizRepo.Delete")
	defer span.End()

	err := repo.db.QueryRowxContext(ctx, "DELETE FROM quizzes WHERE id = $1", id).Err()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}
