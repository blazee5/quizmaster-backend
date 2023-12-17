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

func (repo *Repository) Create(ctx context.Context, input domain.SignUpRequest) (int, error) {
	ctx, span := repo.tracer.Start(ctx, "admin.userRepo.Create")
	defer span.End()

	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		input.Username, input.Email, input.Password).Scan(&id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetAll(ctx context.Context) ([]models.ShortUser, error) {
	ctx, span := repo.tracer.Start(ctx, "admin.userRepo.GetAll")
	defer span.End()

	users := make([]models.ShortUser, 0)

	err := repo.db.SelectContext(ctx, &users, "SELECT id, username, email, avatar FROM users")

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return users, nil

}

func (repo *Repository) Update(ctx context.Context, id int, input domain.User) error {
	ctx, span := repo.tracer.Start(ctx, "admin.userRepo.Update")
	defer span.End()

	err := repo.db.QueryRowxContext(ctx, "UPDATE users SET username = $1, email = $2 WHERE id = $3",
		input.Username, input.Email, id).Err()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, id int) error {
	ctx, span := repo.tracer.Start(ctx, "admin.userRepo.Delete")
	defer span.End()

	err := repo.db.QueryRowxContext(ctx, "DELETE FROM users WHERE id = $1", id).Err()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}
