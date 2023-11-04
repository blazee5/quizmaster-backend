package repository

import (
	"context"
	"database/sql"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
	"slices"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (repo *Repository) GetById(ctx context.Context, userId int) (models.User, error) {
	var user models.User

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE id = $1", userId).StructScan(&user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repo *Repository) GetQuizzes(ctx context.Context, userId int) ([]models.Quiz, error) {
	var quizzes []models.Quiz

	err := repo.db.SelectContext(ctx, &quizzes, "SELECT * FROM quizzes WHERE user_id = $1", userId)

	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (repo *Repository) GetResults(ctx context.Context, userId int) ([]models.Quiz, error) {
	var quizzes []models.Quiz

	rows, err := repo.db.QueryxContext(ctx, "SELECT quiz_id FROM results WHERE user_id = $1", userId)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var quizId int
		var quiz models.Quiz

		err := rows.Scan(&quizId)

		if err != nil {
			return nil, err
		}

		err = repo.db.QueryRowxContext(ctx, "SELECT * FROM quizzes WHERE id = $1", quizId).StructScan(&quiz)

		if err != nil {
			return nil, err
		}

		if !slices.Contains(quizzes, quiz) {
			quizzes = append(quizzes, quiz)
		}
	}

	return quizzes, nil

}

func (repo *Repository) ChangeAvatar(ctx context.Context, userId int, file string) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE users SET avatar = $1 WHERE id = $2", file, userId)

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Update(ctx context.Context, userId int, input domain.UpdateUser) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE users SET fio = COALESCE(NULLIF($1, ''), fio) WHERE id = $2", input.Fio, userId)

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, userId int) error {
	res, err := repo.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userId)

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
