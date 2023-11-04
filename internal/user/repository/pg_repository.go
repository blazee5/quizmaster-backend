package repository

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) GetUserById(ctx context.Context, userId int) (models.User, error) {
	var user models.User

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE id = $1", userId).StructScan(&user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
