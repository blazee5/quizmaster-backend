package repository

import (
	"context"
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"math"
)

type Repository struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewRepository(db *sqlx.DB, tracer trace.Tracer) *Repository {
	return &Repository{db: db, tracer: tracer}
}

func (repo *Repository) GetAll(ctx context.Context, sortBy, sortDir string, page, size int) (models.QuizList, error) {
	ctx, span := repo.tracer.Start(ctx, "quizRepo.GetAll")
	defer span.End()

	var total int
	var offset int

	if size > 0 {
		offset = (page - 1) * size
	}

	err := repo.db.QueryRowxContext(ctx, "SELECT COUNT(*) FROM quizzes").Scan(&total)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.QuizList{}, err
	}

	quizzes := make([]models.Quiz, 0)

	sql, args, err := sq.
		Select("id", "title", "description", "image", "user_id", "created_at").
		From("quizzes").
		OrderBy(sortBy + " " + sortDir).
		Limit(uint64(size)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.QuizList{}, err
	}

	err = repo.db.SelectContext(ctx, &quizzes, sql, args...)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.QuizList{}, err
	}

	return models.QuizList{
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(size))),
		Page:       page,
		Size:       size,
		Quizzes:    quizzes,
	}, nil
}

func (repo *Repository) Create(ctx context.Context, userID int, input domain.Quiz) (models.Quiz, error) {
	ctx, span := repo.tracer.Start(ctx, "quizRepo.Create")
	defer span.End()

	var quiz models.Quiz

	err := repo.db.QueryRowxContext(ctx, `INSERT INTO quizzes (title, description, user_id)
		VALUES ($1, $2, $3) RETURNING id, title, description, image, user_id, created_at`,
		input.Title, input.Description, userID).StructScan(&quiz)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.Quiz{}, err
	}

	return quiz, nil
}

func (repo *Repository) GetByID(ctx context.Context, id int) (models.Quiz, error) {
	ctx, span := repo.tracer.Start(ctx, "quizRepo.GetByID")
	defer span.End()

	var quiz models.Quiz

	err := repo.db.QueryRowxContext(ctx, "SELECT * FROM quizzes WHERE id = $1", id).StructScan(&quiz)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.Quiz{}, err
	}

	return quiz, nil
}

func (repo *Repository) Update(ctx context.Context, quizID int, input domain.Quiz) (models.Quiz, error) {
	ctx, span := repo.tracer.Start(ctx, "quizRepo.Update")
	defer span.End()

	var quiz models.Quiz

	err := repo.db.QueryRowxContext(ctx, `UPDATE quizzes SET
		title = COALESCE(NULLIF($1, ''), title),
		description = COALESCE(NULLIF($2, ''), description) WHERE id = $3
		RETURNING id, title, description, image, user_id, created_at`,
		input.Title, input.Description, quizID).StructScan(&quiz)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return models.Quiz{}, err
	}

	return quiz, nil
}

func (repo *Repository) Delete(ctx context.Context, id int) error {
	ctx, span := repo.tracer.Start(ctx, "quizRepo.Delete")
	defer span.End()

	res, err := repo.db.ExecContext(ctx, "DELETE FROM quizzes WHERE id = $1", id)

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

func (repo *Repository) UploadImage(ctx context.Context, id int, filename string) error {
	ctx, span := repo.tracer.Start(ctx, "quizRepo.UploadImage")
	defer span.End()

	err := repo.db.QueryRowxContext(ctx, "UPDATE quizzes SET image = $1 WHERE id = $2", filename, id).Err()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil
}

func (repo *Repository) DeleteImage(ctx context.Context, id int) error {
	ctx, span := repo.tracer.Start(ctx, "quizRepo.DeleteImage")
	defer span.End()

	err := repo.db.QueryRowxContext(ctx, "UPDATE quizzes SET image = '' WHERE id = $1", id).Err()

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return err
	}

	return nil

}
