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

func (repo *Repository) Create(ctx context.Context, userID int, input domain.Quiz) (int, error) {
	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO quizzes (title, description, user_id) VALUES ($1, $2, $3) RETURNING id",
		input.Title, input.Description, userID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetQuizzes(ctx context.Context) ([]models.Quiz, error) {
	quizzes := make([]models.Quiz, 0)

	err := repo.db.SelectContext(ctx, &quizzes, "SELECT id, title, description, image, user_id, created_at FROM quizzes")

	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (repo *Repository) Update(ctx context.Context, id int, input domain.Quiz) error {
	err := repo.db.QueryRowxContext(ctx, "UPDATE quizzes SET title = $1, description = $2 WHERE id = $3",
		input.Title, input.Description, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, id int) error {
	err := repo.db.QueryRowxContext(ctx, "DELETE FROM quizzes WHERE id = $1", id).Err()

	if err != nil {
		return err
	}

	return nil
}
