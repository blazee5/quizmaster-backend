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

func (repo *Repository) Create(ctx context.Context, userID int, input domain.Quiz) (int, error) {
	var quizID int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO quizzes (title, description, user_id) VALUES ($1, $2, $3) RETURNING id",
		input.Title, input.Description, userID).Scan(&quizID)

	if err != nil {
		return 0, err
	}

	return quizID, nil
}

func (repo *Repository) GetByID(ctx context.Context, id int) (models.Quiz, error) {
	var quiz models.Quiz

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM quizzes WHERE id = $1", id).StructScan(&quiz)

	if err != nil {
		return models.Quiz{}, err
	}

	return quiz, nil
}

func (repo *Repository) Update(ctx context.Context, quizID int, input domain.Quiz) error {
	err := repo.db.QueryRowxContext(ctx, `UPDATE quizzes SET
		title = COALESCE(NULLIF($1, ''), title),
		description = COALESCE(NULLIF($2, ''), description) WHERE id = $3`,
		input.Title, input.Description, quizID).Err()

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

func (repo *Repository) GetAnswerByID(ctx context.Context, id int) (models.Answer, error) {
	var answer models.Answer

	err := repo.db.QueryRowxContext(ctx, selectAnswerQuery, id).StructScan(&answer)
	if err != nil {
		return models.Answer{}, err
	}
	return answer, nil
}

func (repo *Repository) GetAnswersByID(ctx context.Context, id int) ([]models.Answer, error) {
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

func (repo *Repository) SaveUserAnswer(ctx context.Context, tx *sqlx.Tx, userID, questionID, answerID int, answerText string) error {
	_, err := tx.ExecContext(ctx, insertUserAnswerQuery, userID, questionID, answerID, answerText)

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

func (repo *Repository) SaveResult(ctx context.Context, userID, quizID int, score, percent int) error {
	_, err := repo.db.ExecContext(ctx, insertResultQuery, userID, quizID, score, percent)

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
	err := repo.db.QueryRowxContext(ctx, "UPDATE quizzes SET image = '' WHERE id = $1", id).Err()

	if err != nil {
		return err
	}

	return nil

}
