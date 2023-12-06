package repository

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewRepository(db *sqlx.DB, tracer trace.Tracer) *Repository {
	return &Repository{db: db, tracer: tracer}
}

func (repo *Repository) Create(ctx context.Context, questionID int) (int, error) {
	ctx, span := repo.tracer.Start(ctx, "answerRepo.Create")
	defer span.End()

	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO answers (question_id) VALUES ($1) RETURNING id",
		questionID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetByID(ctx context.Context, ID int) (models.Answer, error) {
	ctx, span := repo.tracer.Start(ctx, "answerRepo.GetByID")
	defer span.End()

	var answer models.Answer

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM answers WHERE id = $1", ID).StructScan(&answer)

	if err != nil {
		return models.Answer{}, err
	}

	return answer, nil
}

func (repo *Repository) GetAnswersByQuestionID(ctx context.Context, questionID int) ([]models.Answer, error) {
	ctx, span := repo.tracer.Start(ctx, "answerRepo.GetAnswersByQuestionID")
	defer span.End()

	answers := make([]models.Answer, 0)

	err := repo.db.SelectContext(ctx, &answers, "SELECT * FROM answers WHERE question_id = $1", questionID)

	if err != nil {
		return nil, err
	}

	return answers, nil
}

func (repo *Repository) GetAnswersInfoByQuestionID(ctx context.Context, questionID int) ([]models.AnswerInfo, error) {
	ctx, span := repo.tracer.Start(ctx, "answerRepo.GetAnswersInfoByQuestionID")
	defer span.End()

	answers := make([]models.AnswerInfo, 0)

	err := repo.db.SelectContext(ctx, &answers, `SELECT a.id, a.text, a.question_id, a.order_id FROM answers a
		WHERE a.question_id = $1 ORDER BY a.order_id ASC`, questionID)

	if err != nil {
		return nil, err
	}

	return answers, nil
}

func (repo *Repository) Update(ctx context.Context, answerID int, input domain.Answer) error {
	ctx, span := repo.tracer.Start(ctx, "answerRepo.Update")
	defer span.End()

	err := repo.db.QueryRowxContext(ctx, `UPDATE answers SET
		text = $1,
		is_correct = $2
		WHERE id = $3`,
		input.Text, input.IsCorrect, answerID).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, answerID int) error {
	ctx, span := repo.tracer.Start(ctx, "answerRepo.Delete")
	defer span.End()

	err := repo.db.QueryRowxContext(ctx, `DELETE FROM answers WHERE id = $1`, answerID).Err()

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) ChangeOrder(ctx context.Context, questionID int, input domain.ChangeAnswerOrder) error {
	ctx, span := repo.tracer.Start(ctx, "answerRepo.ChangeOrder")
	defer span.End()

	tx, err := repo.db.Beginx()

	if err != nil {
		return err
	}

	for _, item := range input.Orders {
		_, err := tx.ExecContext(ctx, "UPDATE answers SET order_id = $1 WHERE id = $2 AND question_id = $3", item.OrderID, item.AnswerID, questionID)

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
