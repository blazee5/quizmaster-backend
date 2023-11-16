package repository

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Create(ctx context.Context, questionId int) (int, error) {
	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO answers (question_id) VALUES ($1) RETURNING id",
		questionId).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) Update(ctx context.Context, answerId int, input domain.Answer) error {
	err := repo.db.QueryRowxContext(ctx, `UPDATE answers SET
		text = COALESCE(NULLIF($1, ''), text),
		is_correct = $2
		WHERE id = $3`,
		input.Text, input.IsCorrect, answerId).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, answerId int) error {
	err := repo.db.QueryRowxContext(ctx, `DELETE FROM answers WHERE id = $1`, answerId).Err()

	if err != nil {
		return err
	}

	return nil
}
