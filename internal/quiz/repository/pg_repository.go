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
