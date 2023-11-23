package repository

import (
	"context"
	"database/sql"
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

func (repo *Repository) CreateQuestion(ctx context.Context, quizID int) (int, error) {
	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO questions (quiz_id) VALUES ($1) RETURNING id",
		quizID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetQuestionsByID(ctx context.Context, id int) ([]models.Question, error) {
	questions := make([]models.Question, 0)

	if err := repo.db.SelectContext(ctx, &questions, "SELECT * FROM questions WHERE quiz_id = $1", id); err != nil {
		return nil, err
	}

	if len(questions) == 0 {
		return nil, sql.ErrNoRows
	}

	answers := make([]models.Answer, 0)

	query := `
        SELECT id, text, question_id, order_id
        FROM answers
        WHERE question_id IN (
            SELECT id
            FROM questions
            WHERE quiz_id = $1
        )`

	if err := repo.db.SelectContext(ctx, &answers, query, id); err != nil {
		return nil, err
	}

	for i := range questions {
		for _, answer := range answers {
			if answer.QuestionID == questions[i].ID {
				questions[i].Answers = append(questions[i].Answers, answer)
			}
		}
	}

	return questions, nil
}

func (repo *Repository) GetQuestionsWithAnswers(ctx context.Context, id int) ([]models.QuestionWithAnswers, error) {
	questions := make([]models.QuestionWithAnswers, 0)

	if err := repo.db.SelectContext(ctx, &questions, "SELECT * FROM questions WHERE quiz_id = $1", id); err != nil {
		return nil, err
	}

	if len(questions) == 0 {
		return nil, sql.ErrNoRows
	}

	answers := make([]models.AnswerInfo, 0)

	query := `
        SELECT id, text, question_id, order_id, is_correct
        FROM answers
        WHERE question_id IN (
            SELECT id
            FROM questions
            WHERE quiz_id = $1
        )`

	if err := repo.db.SelectContext(ctx, &answers, query, id); err != nil {
		return nil, err
	}

	for i := range questions {
		for _, answer := range answers {
			if answer.QuestionID == questions[i].ID {
				questions[i].Answers = append(questions[i].Answers, answer)
			}
		}
	}

	return questions, nil
}

func (repo *Repository) Update(ctx context.Context, id int, input domain.Question) error {
	err := repo.db.QueryRowxContext(ctx, `UPDATE questions
		SET title = $1,
		    type = COALESCE(NULLIF($2, ''), type)
		WHERE id = $3`,
		input.Title, input.Type, id).Err()

	if err != nil {
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

func (repo *Repository) UploadImage(ctx context.Context, id int, filename string) error {
	err := repo.db.QueryRowxContext(ctx, "UPDATE questions SET image = $1 WHERE id = $2", filename, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) DeleteImage(ctx context.Context, id int) error {
	err := repo.db.QueryRowxContext(ctx, "UPDATE questions SET image = '' WHERE id = $1", id).Err()

	if err != nil {
		return err
	}

	return nil

}
