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
	GetQuestionsById(ctx context.Context, id int, includeIsCorrect bool) ([]models.Question, error)
	GetCorrectAnswers(ctx context.Context, id int) (int, error)
	GetAnswerById(ctx context.Context, id int) (models.Answer, error)
	GetAnswersById(ctx context.Context, id int) ([]models.Answer, error)
	GetQuestionType(ctx context.Context, id int) (string, error)
	Create(ctx context.Context, input domain.Quiz) (int, error)
	SaveResult(ctx context.Context, userId, quizId int, score, percent int) error
	SaveUserAnswer(ctx context.Context, tx *sqlx.Tx, userId, questionId, answerId int, answerText string) error
	Delete(ctx context.Context, quizId int) error
	NewTx() (*sqlx.Tx, error)
}
