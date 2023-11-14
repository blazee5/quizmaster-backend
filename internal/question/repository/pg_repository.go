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
	questions := make([]models.Question, 0)

	if err := repo.db.SelectContext(ctx, &questions, "SELECT * FROM questions WHERE quiz_id = $1", id); err != nil {
		return nil, err
	}

	answers := make([]models.Answer, 0)

	query := `
        SELECT id, text, question_id
        FROM answers
        WHERE question_id IN (
            SELECT id
            FROM questions
            WHERE quiz_id = $1
        )`
	if includeIsCorrect {
		query = `
        SELECT id, text, question_id, is_correct
        FROM answers
        WHERE question_id IN (
            SELECT id
            FROM questions
            WHERE quiz_id = $1
        )`
	}

	if err := repo.db.SelectContext(ctx, &answers, query, id); err != nil {
		return nil, err
	}

	for i := range questions {
		for _, answer := range answers {
			if answer.QuestionId == questions[i].Id {
				questions[i].Answers = append(questions[i].Answers, answer)
			}
		}
	}

	return questions, nil
}

func (repo *Repository) Update(ctx context.Context, id int, input domain.Question) error {
	tx, err := repo.db.Beginx()

	if err != nil {
		return err
	}

	err = tx.QueryRowxContext(ctx, `UPDATE questions
		SET title = COALESCE(NULLIF($1, ''), title),
		    image = COALESCE(NULLIF($2, ''), image),
		    type = COALESCE(NULLIF($3, ''))
		WHERE id = $4`,
		input.Title, input.Image, input.Type, id).Err()

	if err != nil {
		tx.Rollback()
		return err
	}

	for _, answer := range input.Answers {
		_, err = tx.ExecContext(ctx, "UPDATE answers SET text = $1, is_correct = $3 WHERE id = $4",
			answer.Text, answer.IsCorrect, id)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, id int) error {
	err := repo.db.QueryRowxContext(ctx, "DELETE FROM questions WHERE id = $1", id).Err()

	if err != nil {
		return err
	}

	return nil
}
