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

func (repo *Repository) NewTx() (*sqlx.Tx, error) {
	tx, err := repo.db.Beginx()

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (repo *Repository) GetAll(ctx context.Context) ([]models.Quiz, error) {
	quizzes := make([]models.Quiz, 0)

	err := repo.db.SelectContext(ctx, &quizzes, "SELECT * FROM quizzes")

	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (repo *Repository) Create(ctx context.Context, userId int, input domain.Quiz) (int, error) {
	var quizId int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO quizzes (title, description, user_id) VALUES ($1, $2, $3) RETURNING id",
		input.Title, input.Description, userId).Scan(&quizId)

	if err != nil {
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

func (repo *Repository) Update(ctx context.Context, quizId int, input domain.Quiz) error {
	err := repo.db.QueryRowxContext(ctx, `UPDATE quizzes SET
		title = COALESCE(NULLIF($1, ''), title),
		description = COALESCE(NULLIF($2, ''), description) WHERE id = $3`,
		input.Title, input.Description, quizId).Err()

	if err != nil {
		return err
	}

	return nil
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

func (repo *Repository) GetAnswerById(ctx context.Context, id int) (models.Answer, error) {
	var answer models.Answer

	err := repo.db.QueryRowxContext(ctx, selectAnswerQuery, id).StructScan(&answer)
	if err != nil {
		return models.Answer{}, err
	}
	return answer, nil
}

func (repo *Repository) GetAnswersById(ctx context.Context, id int) ([]models.Answer, error) {
	var answer []models.Answer

	err := repo.db.SelectContext(ctx, &answer, selectAnswersQuery, id)

	if err != nil {
		return nil, err
	}

	return answer, nil
}

func (repo *Repository) GetQuestionType(ctx context.Context, id int) (string, error) {
	var questionType string

	err := repo.db.QueryRowContext(ctx, selectQuestionTypeQuery, id).Scan(&questionType)
	if err != nil {
		return "", err
	}

	return questionType, nil
}

func (repo *Repository) SaveUserAnswer(ctx context.Context, tx *sqlx.Tx, userId, questionId, answerId int, answerText string) error {
	_, err := tx.ExecContext(ctx, insertUserAnswerQuery, userId, questionId, answerId, answerText)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (repo *Repository) GetCorrectAnswers(ctx context.Context, id int) (int, error) {
	var totalCorrectAnswers int

	err := repo.db.QueryRowContext(ctx, selectTotalCorrectQuery, id).Scan(&totalCorrectAnswers)

	if err != nil {
		return 0, err
	}

	return totalCorrectAnswers, nil
}

func (repo *Repository) SaveResult(ctx context.Context, userId, quizId int, score, percent int) error {
	_, err := repo.db.ExecContext(ctx, insertResultQuery, userId, quizId, score, percent)

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) UploadImage(ctx context.Context, id int, filename string) error {
	err := repo.db.QueryRowxContext(ctx, "UPDATE quizzes SET image = $1 WHERE id = $2", filename, id).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) DeleteImage(ctx context.Context, id int) error {
	err := repo.db.QueryRowxContext(ctx, "UPDATE quizzes SET image = null WHERE id = $1", id).Err()

	if err != nil {
		return err
	}

	return nil

}
