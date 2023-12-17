package repository

import (
	"context"
	"database/sql"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"slices"
)

type Repository struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewRepository(db *sqlx.DB, tracer trace.Tracer) *Repository {
	return &Repository{db: db, tracer: tracer}
}

func (repo *Repository) GetByID(ctx context.Context, userID int) (models.UserInfo, error) {
	ctx, span := repo.tracer.Start(ctx, "userRepo.GetByID")
	defer span.End()

	var user models.ShortUser

	err := repo.db.QueryRowxContext(ctx, "SELECT id, username, email, avatar FROM users WHERE id = $1", userID).StructScan(&user)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.UserInfo{}, err
	}

	quizzes := make([]models.Quiz, 0)

	err = repo.db.SelectContext(ctx, &quizzes, "SELECT * FROM quizzes WHERE user_id = $1", userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.UserInfo{}, err
	}

	userResults := make([]models.UserResult, 0)
	processedQuizzes := make([]int, 0)

	query := `SELECT q.id, q.title, q.description, q.image, q.user_id, q.created_at, r.score,
       (SELECT COUNT(*) FROM questions WHERE questions.quiz_id = q.id) AS questions_count, r.created_at
		FROM results r
		INNER JOIN quizzes q ON r.quiz_id = q.id
		WHERE r.user_id = $1 AND r.is_completed = true
		GROUP BY q.id, r.score, r.is_completed, r.created_at
		ORDER BY r.score DESC`

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.UserInfo{}, err
	}

	rows, err := repo.db.QueryxContext(ctx, query, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.UserInfo{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var userResult models.UserResult
		var quiz models.Quiz

		err = rows.Scan(
			&quiz.ID,
			&quiz.Title,
			&quiz.Description,
			&quiz.Image,
			&quiz.UserID,
			&quiz.CreatedAt,
			&userResult.Score,
			&userResult.QuestionsCount,
			&userResult.CreatedAt,
		)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return models.UserInfo{}, err
		}

		if !slices.Contains(processedQuizzes, quiz.ID) {
			userResult.Quiz = quiz
			userResults = append(userResults, userResult)
			processedQuizzes = append(processedQuizzes, quiz.ID)
		}
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.UserInfo{}, err
	}

	return models.UserInfo{
		User:    user,
		Quizzes: quizzes,
		Results: userResults,
	}, nil
}

func (repo *Repository) GetQuizzes(ctx context.Context, userID int) ([]models.Quiz, error) {
	ctx, span := repo.tracer.Start(ctx, "userRepo.GetQuizzes")
	defer span.End()

	quizzes := make([]models.Quiz, 0)

	err := repo.db.SelectContext(ctx, &quizzes, "SELECT * FROM quizzes WHERE user_id = $1", userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	return quizzes, nil
}

func (repo *Repository) GetResults(ctx context.Context, userID int) ([]models.Quiz, error) {
	ctx, span := repo.tracer.Start(ctx, "userRepo.GetResults")
	defer span.End()

	quizzes := make([]models.Quiz, 0)

	rows, err := repo.db.QueryxContext(ctx, "SELECT quiz_id FROM results WHERE user_id = $1", userID)
	defer rows.Close()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	for rows.Next() {
		var quizID int
		var quiz models.Quiz

		err := rows.Scan(&quizID)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return nil, err
		}

		err = repo.db.QueryRowxContext(ctx, "SELECT * FROM quizzes WHERE id = $1", quizID).StructScan(&quiz)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			return nil, err
		}

		if !slices.Contains(quizzes, quiz) {
			quizzes = append(quizzes, quiz)
		}
	}

	return quizzes, nil
}

func (repo *Repository) GetAvatarByID(ctx context.Context, userID int) (string, error) {
	ctx, span := repo.tracer.Start(ctx, "userRepo.GetAvatarByID")
	defer span.End()

	var avatar string

	err := repo.db.QueryRowxContext(ctx, "SELECT avatar FROM users WHERE id = $1", userID).Scan(&avatar)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return "", err
	}

	return avatar, nil
}

func (repo *Repository) ChangeAvatar(ctx context.Context, userID int, file string) error {
	ctx, span := repo.tracer.Start(ctx, "userRepo.ChangeAvatar")
	defer span.End()

	_, err := repo.db.ExecContext(ctx, "UPDATE users SET avatar = $1 WHERE id = $2", file, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (repo *Repository) Update(ctx context.Context, userID int, input domain.UpdateUser) error {
	ctx, span := repo.tracer.Start(ctx, "userRepo.Update")
	defer span.End()

	_, err := repo.db.ExecContext(ctx, "UPDATE users SET username = COALESCE(NULLIF($1, ''), username) WHERE id = $2", input.Username, userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (repo *Repository) Delete(ctx context.Context, userID int) error {
	ctx, span := repo.tracer.Start(ctx, "userRepo.Delete")
	defer span.End()

	res, err := repo.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	if rows < 1 {
		return sql.ErrNoRows
	}

	return nil
}
