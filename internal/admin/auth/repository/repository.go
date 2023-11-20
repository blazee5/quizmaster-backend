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

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) ValidateUser(ctx context.Context, input domain.SignInRequest) (models.User, error) {
	var user models.User

	err := repo.db.QueryRowxContext(ctx, "SELECT users.id, users.role_id FROM users JOIN roles r on r.id = users.role_id WHERE email = $1 AND password = $2 AND r.name = 'admin'",
		input.Email, input.Password).StructScan(&user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
