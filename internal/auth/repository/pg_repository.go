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

func (repo *Repository) UpdateEmail(ctx context.Context, userID int, email string) error {
	err := repo.db.QueryRowxContext(ctx, "UPDATE users SET email = $1 WHERE id = $2", email, userID).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) UpdatePassword(ctx context.Context, userID int, password string) error {
	err := repo.db.QueryRowxContext(ctx, "UPDATE users SET password = $1 WHERE id = $2", password, userID).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) CreateVerificationCode(ctx context.Context, userID int, codeType, code, email string) error {
	ctx, span := repo.tracer.Start(ctx, "authRepo.CreateVerificationCode")
	defer span.End()

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO verification_codes (type, code, user_id, email) VALUES ($1, $2, $3, $4)", codeType, code, userID, email).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetVerificationCode(ctx context.Context, code, codeType string) (models.VerificationCode, error) {
	ctx, span := repo.tracer.Start(ctx, "authRepo.GetVerificationCode")
	defer span.End()

	var verificationCode models.VerificationCode

	err := repo.db.QueryRowxContext(ctx, "SELECT id, type, code, user_id, email, expire_date FROM verification_codes WHERE code = $1 AND type = $2", code, codeType).
		StructScan(&verificationCode)

	if err != nil {
		return models.VerificationCode{}, err
	}

	return verificationCode, nil
}

func (repo *Repository) DeleteVerificationCode(ctx context.Context, id int) error {
	ctx, span := repo.tracer.Start(ctx, "authRepo.DeleteVerificationCode")
	defer span.End()

	err := repo.db.QueryRowxContext(ctx, "DELETE FROM verification_codes WHERE id = $1", id).Err()

	if err != nil {
		return err
	}

	return nil
}
