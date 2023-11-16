package quiz

import (
	"context"
	"github.com/blazee5/quizmaster-backend/internal/domain"
	"github.com/blazee5/quizmaster-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetById(ctx context.Context, id int) (models.Quiz, error)
	GetAll(ctx context.Context) ([]models.Quiz, error)
	GetCorrectAnswers(ctx context.Context, id int) (int, error)
	Create(ctx context.Context, userId int, input domain.Quiz) (int, error)
	SaveResult(ctx context.Context, userId, quizId int, score, percent int) error
	SaveUserAnswer(ctx context.Context, tx *sqlx.Tx, userId, questionId, answerId int, answerText string) error
	Delete(ctx context.Context, quizId int) error
	GetQuestionsById(ctx context.Context, id int, includeIsCorrect bool) ([]models.Question, error)
	GetAnswerById(ctx context.Context, id int) (models.Answer, error)
	GetAnswersById(ctx context.Context, id int) ([]models.Answer, error)
	GetQuestionType(ctx context.Context, id int) (string, error)
	NewTx() (*sqlx.Tx, error)
	UploadImage(ctx context.Context, id int, filename string) error
	DeleteImage(ctx context.Context, id int) error
}
