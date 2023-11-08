package repository

const (
	insertUserAnswerQuery   = `INSERT INTO user_answers (user_id, question_id, answer_id, text) VALUES ($1, $2, $3, $4)`
	insertResultQuery       = `INSERT INTO results (user_id, quiz_id, score, percent) VALUES ($1, $2, $3, $4)`
	selectAnswerQuery       = `SELECT * FROM answers WHERE id = $1`
	selectAnswersQuery      = `SELECT * FROM answers WHERE question_id = $1`
	selectTotalCorrectQuery = `SELECT COUNT(is_correct) FROM answers WHERE question_id = $1 AND is_correct = true`
	selectQuestionTypeQuery = `SELECT type FROM questions WHERE id = $1`
)
