package repository

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return Repository{db: db}
}

func (repo *Repository) CreateUser(ctx context.Context, input domain.SignUpRequest) (int, error) {
	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		input.Username, input.Email, input.Password).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) ValidateUser(ctx context.Context, input domain.SignInRequest) (models.User, error) {
	var user models.User

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE email = $1 AND password = $2", input.Email, input.Password).StructScan(&user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
