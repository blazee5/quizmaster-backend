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

func (repo *Repository) GetAll(ctx context.Context) ([]models.Quiz, error) {
	quizzes := make([]models.Quiz, 0)

	err := repo.db.SelectContext(ctx, &quizzes, "SELECT * FROM quizzes")

	if err != nil {
		return nil, err
	}

	return quizzes, nil
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

func (repo *Repository) SaveResult(ctx context.Context, userId int, quizId int, input domain.Result) (int, error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var score float64
	var totalQuestions float64

	err = tx.QueryRowContext(ctx, selectTotalQuestionsQuery, quizId).Scan(&totalQuestions)

	if err != nil {
		return 0, err
	}

	for questionId, answerIds := range input.Answers {
		var totalCorrectAnswers int
		var userCorrectAnswers int
		var totalUserAnswers int

		err := tx.QueryRowContext(ctx, selectTotalCorrectQuery, questionId).Scan(&totalCorrectAnswers)

		if err != nil {
			return 0, err
		}

		for _, answerId := range answerIds {
			var isCorrect bool
			err := tx.QueryRowContext(ctx, selectAnswerQuery, answerId).Scan(&isCorrect)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			_, err = tx.ExecContext(ctx, insertUserAnswerQuery, userId, questionId, answerId, isCorrect)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			if isCorrect {
				userCorrectAnswers++
			}

			totalUserAnswers++
		}

		if userCorrectAnswers == totalCorrectAnswers && totalUserAnswers == userCorrectAnswers {
			score++
		}
	}

	_, err = tx.ExecContext(ctx, insertResultQuery, userId, quizId, score, score/totalQuestions*100)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(score), nil
}

func (repo *Repository) Delete(ctx context.Context, id int) error {
	res, err := repo.db.ExecContext(ctx, "DELETE FROM quizzes WHERE id = $1", id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rows < 1 {
		return sql.ErrNoRows
	}

	return nil
}
