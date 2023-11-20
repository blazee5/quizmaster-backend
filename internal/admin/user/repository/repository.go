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

func (repo *Repository) Create(ctx context.Context, input domain.SignUpRequest) (int, error) {
	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		input.Username, input.Email, input.Password).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetAll(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0)

	err := repo.db.SelectContext(ctx, &users, "SELECT * FROM users")

	if err != nil {
		return make([]models.User, 0), err
	}

	return users, nil

}

func (repo *Repository) Update(ctx context.Context, id int, input domain.User) error {
	err := repo.db.QueryRowxContext(ctx, "UPDATE users SET username = $1, email = $2 WHERE id = $3",
		input.Username, input.Email, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, id int) error {
	err := repo.db.QueryRowxContext(ctx, "DELETE FROM users WHERE id = $1", id).Err()

	if err != nil {
		return err
	}

	return nil
}
