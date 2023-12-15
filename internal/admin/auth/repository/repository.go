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

func (repo *Repository) ValidateUser(ctx context.Context, input domain.SignInRequest) (models.User, error) {
	ctx, span := repo.tracer.Start(ctx, "admin.authRepo.ValidateUser")
	defer span.End()

	var user models.User

	err := repo.db.QueryRowxContext(ctx, "SELECT users.id, users.role_id FROM users JOIN roles r on r.id = users.role_id WHERE email = $1 AND password = $2 AND r.name = 'admin'",
		input.Email, input.Password).StructScan(&user)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.User{}, err
	}

	return user, nil
}
