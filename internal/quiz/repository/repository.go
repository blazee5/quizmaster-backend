package repository

import (
	"context"
	"github.com/blazee5/testhub-backend/internal/domain"
	"github.com/blazee5/testhub-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) Create(ctx context.Context, input domain.Quiz) (int, error) {
	var quizId int

	tx, err := repo.db.Beginx()

	if err != nil {
		return 0, err
	}

	err = tx.QueryRowxContext(ctx, "INSERT INTO quizzes (title, description, image, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
		input.Title, input.Description, input.Image, input.UserId).Scan(&quizId)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, question := range input.Questions {
		var questionId int
		question.QuizId = quizId

		err := tx.QueryRowxContext(ctx, "INSERT INTO questions (title, image, quiz_id) VALUES ($1, $2, $3) RETURNING id",
			question.Title, question.Image, quizId).Scan(&questionId)

		if err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, answer := range question.Answers {
			answer.QuestionId = questionId

			_, err := tx.ExecContext(ctx, "INSERT INTO answers (text, is_correct, question_id) VALUES ($1, $2, $3)",
				answer.Text, answer.IsCorrect, answer.QuestionId)

			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return quizId, nil
}

func (repo *Repository) GetById(ctx context.Context, id int) (models.Quiz, error) {
	var quiz models.Quiz

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM quizzes WHERE id = $1", id).StructScan(&quiz)

	if err != nil {
		return models.Quiz{}, err
	}

	return quiz, nil
}

func (repo *Repository) GetQuestionsById(ctx context.Context, id int) ([]models.Question, error) {
	var questions []models.Question

	if err := repo.db.SelectContext(ctx, &questions, "SELECT * FROM questions WHERE quiz_id = $1", id); err != nil {
		return nil, err
	}

	var answers []models.Answer
	if err := repo.db.SelectContext(ctx, &answers, `
        SELECT id, text, question_id
        FROM answers
        WHERE question_id IN (
            SELECT id
            FROM questions
            WHERE quiz_id = $1
        )`, id); err != nil {
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
