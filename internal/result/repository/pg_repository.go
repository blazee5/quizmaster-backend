package repository

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Repository struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewRepository(db *sqlx.DB, tracer trace.Tracer) *Repository {
	return &Repository{db: db, tracer: tracer}
}

func (repo *Repository) NewResult(ctx context.Context, userID, quizID int) (int, error) {
	ctx, span := repo.tracer.Start(ctx, "resultRepo.NewResult")
	defer span.End()

	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO results (user_id, quiz_id, score) VALUES ($1, $2, $3) RETURNING id",
		userID, quizID, 0).Scan(&id)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetByID(ctx context.Context, id int) (models.Result, error) {
	ctx, span := repo.tracer.Start(ctx, "resultRepo.GetByID")
	defer span.End()

	var result models.Result

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM results WHERE id = $1", id).StructScan(&result)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.Result{}, err
	}

	return result, nil
}

func (repo *Repository) GetByQuizID(ctx context.Context, quizID int) ([]models.UsersResult, error) {
	ctx, span := repo.tracer.Start(ctx, "resultRepo.GetByQuizID")
	defer span.End()

	var results []models.UsersResult

	err := repo.db.SelectContext(ctx, &results, `SELECT r.id, r.score, r.created_at, u.username FROM results r
	INNER JOIN (
    	SELECT user_id, MAX(score) AS best_score
    	FROM results
    	WHERE quiz_id = $1 AND is_completed = true
    	GROUP BY user_id
    ) AS best_results ON best_results.user_id = r.user_id AND best_results.best_score = r.score
		JOIN users u ON u.id = r.user_id
		ORDER BY r.score DESC`, quizID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return results, nil
}

func (repo *Repository) GetByUserID(ctx context.Context, id int) (models.Result, error) {
	ctx, span := repo.tracer.Start(ctx, "resultRepo.GetByUserID")
	defer span.End()

	var result models.Result

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM results WHERE user_id = $1 AND is_completed = false", id).StructScan(&result)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.Result{}, err
	}

	return result, nil
}

func (repo *Repository) GetUserAnswerByID(ctx context.Context, answerID, resultID int) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "resultRepo.GetUserAnswerByID")
	defer span.End()

	var result int

	err := repo.db.QueryRowxContext(ctx, "SELECT COUNT(id) FROM user_answers WHERE answer_id = $1 AND result_id = $2", answerID, resultID).Scan(&result)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return false, err
	}

	return result > 0, nil
}

func (repo *Repository) SaveUserAnswer(ctx context.Context, userID, questionID, answerID, resultID int, answerText string) error {
	ctx, span := repo.tracer.Start(ctx, "resultRepo.SaveUserAnswer")
	defer span.End()

	_, err := repo.db.ExecContext(ctx, "INSERT INTO user_answers (user_id, question_id, answer_id, result_id, text) VALUES ($1, $2, $3, $4, $5)",
		userID, questionID, answerID, resultID, answerText)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (repo *Repository) UpdateResult(ctx context.Context, id, userID, score int) error {
	ctx, span := repo.tracer.Start(ctx, "resultRepo.UpdateResult")
	defer span.End()

	_, err := repo.db.ExecContext(ctx, "UPDATE results SET score = score + $1 WHERE id = $2 AND user_id = $3",
		score, id, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (repo *Repository) SubmitResult(ctx context.Context, userID, resultID int) (models.UsersResult, error) {
	ctx, span := repo.tracer.Start(ctx, "resultRepo.SubmitResult")
	defer span.End()

	var result models.UsersResult

	err := repo.db.QueryRowxContext(ctx, "UPDATE results SET is_completed = true WHERE id = $1 AND user_id = $2", resultID, userID).Err()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.UsersResult{}, nil
	}

	err = repo.db.QueryRowxContext(ctx, "SELECT id, score, created_at FROM results WHERE id = $1", resultID).StructScan(&result)

	return result, nil
}
