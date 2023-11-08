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

func (repo *Repository) GetById(ctx context.Context, userId int) (models.UserInfo, error) {
	var user models.ShortUser

	err := repo.db.QueryRowxContext(ctx, "SELECT id, username, email, avatar FROM users WHERE id = $1", userId).StructScan(&user)

	if err != nil {
		return models.UserInfo{}, err
	}

	quizzes := make([]models.Quiz, 0)

	err = repo.db.SelectContext(ctx, &quizzes, "SELECT * FROM quizzes WHERE user_id = $1", userId)

	if err != nil {
		return models.UserInfo{}, err
	}

	var userResults []models.UserResult
	processedQuizzes := make([]int, 0)

	query := `SELECT q.id AS quiz_id, q.title, q.description, q.image, q.user_id, q.created_at, r.score, r.percent, r.created_at
			  FROM results r
		      INNER JOIN quizzes q ON r.quiz_id = q.id 
		      WHERE r.user_id = $1 
			  ORDER BY r.score DESC`

	if err != nil {
		return models.UserInfo{}, err
	}

	rows, err := repo.db.QueryxContext(ctx, query, userId)
	if err != nil {
		return models.UserInfo{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var userResult models.UserResult
		var quiz models.Quiz

		err := rows.Scan(
			&quiz.Id,
			&quiz.Title,
			&quiz.Description,
			&quiz.Image,
			&quiz.UserId,
			&quiz.CreatedAt,
			&userResult.Score,
			&userResult.Percent,
			&userResult.CreatedAt,
		)
		if err != nil {
			return models.UserInfo{}, err
		}

		if !slices.Contains(processedQuizzes, quiz.Id) {
			userResult.Quiz = quiz
			userResults = append(userResults, userResult)
			processedQuizzes = append(processedQuizzes, quiz.Id)
		}
	}

	if err := rows.Err(); err != nil {
		return models.UserInfo{}, err
	}

	return models.UserInfo{
		User:    user,
		Quizzes: quizzes,
		Results: userResults,
	}, nil
}

func (repo *Repository) GetQuizzes(ctx context.Context, userId int) ([]models.Quiz, error) {
	quizzes := make([]models.Quiz, 0)

	err := repo.db.SelectContext(ctx, &quizzes, "SELECT * FROM quizzes WHERE user_id = $1", userId)

	if err != nil {
		return nil, err
	}

	return quizzes, nil
}

func (repo *Repository) GetResults(ctx context.Context, userId int) ([]models.Quiz, error) {
	quizzes := make([]models.Quiz, 0)

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
	_, err := repo.db.ExecContext(ctx, "UPDATE users SET username = COALESCE(NULLIF($1, ''), username) WHERE id = $2", input.Username, userId)

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
