package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetByID(ctx context.Context, id int) (models.Quiz, error)
	GetAll(ctx context.Context) ([]models.Quiz, error)
	GetCorrectAnswers(ctx context.Context, id int) (int, error)
	Create(ctx context.Context, userID int, input domain.Quiz) (int, error)
	SaveResult(ctx context.Context, userID, quizID int, score, percent int) error
	SaveUserAnswer(ctx context.Context, tx *sqlx.Tx, userID, questionID, answerID int, answerText string) error
	Update(ctx context.Context, quizID int, input domain.Quiz) error
	Delete(ctx context.Context, quizID int) error
	GetAnswerByID(ctx context.Context, id int) (models.Answer, error)
	GetAnswersByID(ctx context.Context, id int) ([]models.Answer, error)
	GetQuestionType(ctx context.Context, id int) (string, error)
	NewTx() (*sqlx.Tx, error)
	UploadImage(ctx context.Context, id int, filename string) error
	DeleteImage(ctx context.Context, id int) error
}
