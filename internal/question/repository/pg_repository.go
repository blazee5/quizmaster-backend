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

func (repo *Repository) CreateQuestion(ctx context.Context, quizId int, input domain.Question) (int, error) {
	var id int

	tx, err := repo.db.Beginx()

	if err != nil {
		return 0, err
	}

	err = tx.QueryRowxContext(ctx, "INSERT INTO questions (title, image, quiz_id, type) VALUES ($1, $2, $3, $4) RETURNING id",
		input.Title, input.Image, quizId, input.Type).Scan(&id)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, answer := range input.Answers {
		_, err = tx.ExecContext(ctx, "INSERT INTO answers (text, question_id, is_correct) VALUES ($1, $2, $3)",
			answer.Text, id, answer.IsCorrect)

		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetQuestionsById(ctx context.Context, id int, includeIsCorrect bool) ([]models.Question, error) {
	return nil, nil
}
