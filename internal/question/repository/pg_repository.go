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

func (repo *Repository) CreateQuestion(ctx context.Context, quizID int) (int, error) {
	var id int

	err := repo.db.QueryRowxContext(ctx, "INSERT INTO questions (quiz_id) VALUES ($1) RETURNING id",
		quizID).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) GetQuestionsByID(ctx context.Context, id int) ([]models.Question, error) {
	questions := make([]models.Question, 0)

	rows, err := repo.db.QueryxContext(ctx,
		`SELECT q.id, q.title, q.image, q.quiz_id, q.type, q.order_id, a.id,
     	a.text, a.question_id, a.order_id
		FROM questions q
		JOIN answers a ON a.question_id = q.id
		WHERE quiz_id = $1 AND q.type = 'choice'
		ORDER BY q.order_id, a.order_id ASC
	`, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var question models.Question
		var answer models.Answer

		err = rows.Scan(&question.ID, &question.Title, &question.Image, &question.QuizID, &question.Type, &question.OrderID, &answer.ID, &answer.Text, &answer.QuestionID, &answer.OrderID)

		if err != nil {
			return nil, err
		}

		question.Answers = append(question.Answers, answer)

		if !slices.ContainsFunc(questions, func(n models.Question) bool {
			return n.ID == question.ID
		}) {
			questions = append(questions, question)
		} else {
			idx := slices.IndexFunc(questions, func(n models.Question) bool {
				return n.ID == question.ID
			})

			questions[idx].Answers = append(questions[idx].Answers, answer)
		}

	}

	return questions, nil
}

func (repo *Repository) GetQuestionsWithAnswers(ctx context.Context, id int) ([]models.QuestionWithAnswers, error) {
	questions := make([]models.QuestionWithAnswers, 0)

	if err := repo.db.SelectContext(ctx, &questions, "SELECT * FROM questions WHERE quiz_id = $1 ORDER BY order_id ASC", id); err != nil {
		return nil, err
	}

	if len(questions) == 0 {
		return nil, sql.ErrNoRows
	}

	answers := make([]models.AnswerInfo, 0)

	query := `
        SELECT id, text, question_id, order_id, is_correct
        FROM answers
        WHERE question_id IN (
            SELECT id
            FROM questions
            WHERE quiz_id = $1
        ) ORDER BY order_id ASC`

	if err := repo.db.SelectContext(ctx, &answers, query, id); err != nil {
		return nil, err
	}

	for i := range questions {
		for _, answer := range answers {
			if answer.QuestionID == questions[i].ID {
				questions[i].Answers = append(questions[i].Answers, answer)
			}
		}
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
