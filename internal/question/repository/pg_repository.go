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

func (repo *Repository) CreateQuestion(ctx context.Context, quizID int) (int, error) {
	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO questions (quiz_id) VALUES ($1) RETURNING id",
		quizID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetQuestionByID(ctx context.Context, id int) (models.Question, error) {
	var question models.Question

	if err := repo.db.QueryRowxContext(ctx, "SELECT * FROM questions WHERE id = $1", id).StructScan(&question); err != nil {
		return models.Question{}, err
	}

	return question, nil
}

func (repo *Repository) GetQuestionsByQuizID(ctx context.Context, quizID int) ([]models.Question, error) {
	questions := make([]models.Question, 0)

	err := repo.db.SelectContext(ctx, &questions,
		`SELECT q.id, q.title, q.image, q.quiz_id, q.type, q.order_id FROM questions q
		WHERE quiz_id = $1
		ORDER BY q.order_id ASC`, quizID)

	if err != nil {
		return nil, err
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

func (repo *Repository) ChangeOrder(ctx context.Context, input domain.ChangeQuestionOrder) error {
	tx, err := repo.db.Beginx()

	if err != nil {
		return err
	}

	for _, item := range input.Orders {
		_, err := tx.ExecContext(ctx, "UPDATE questions SET order_id = $1 WHERE id = $2", item.OrderID, item.QuestionID)

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
