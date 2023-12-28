package repository

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewRepository(db *sqlx.DB, tracer trace.Tracer) *Repository {
	return &Repository{db: db, tracer: tracer}
}

func (repo *Repository) CreateUser(ctx context.Context, input domain.SignUpRequest) (int, error) {
	ctx, span := repo.tracer.Start(ctx, "authRepo.CreateUser")
	defer span.End()

	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		input.Username, input.Email, input.Password).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) ValidateUser(ctx context.Context, input domain.SignInRequest) (models.User, error) {
	ctx, span := repo.tracer.Start(ctx, "authRepo.ValidateUser")
	defer span.End()

	var user models.User

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE email = $1 AND password = $2", input.Email, input.Password).StructScan(&user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repo *Repository) CreateVerificationCode(ctx context.Context, userID int, codeType, code string) error {
	err := repo.db.QueryRowxContext(ctx, "INSERT INTO verification_codes (type, code, user_id) VALUES ($1, $2, $3)", codeType, code, userID).Err()

	if err != nil {
		return err
	}

	return nil
}
