package repository

const (
	insertUserAnswerQuery     = `INSERT INTO user_answers (user_id, question_id, answer_id, is_correct) VALUES ($1, $2, $3, $4)`
	insertResultQuery         = `INSERT INTO results (user_id, quiz_id, score, percent) VALUES ($1, $2, $3, $4)`
	selectAnswerQuery         = `SELECT is_correct FROM answers WHERE id = $1`
	selectTotalCorrectQuery   = `SELECT COUNT(is_correct) FROM answers WHERE question_id = $1 AND is_correct = true`
	selectTotalQuestionsQuery = `SELECT COUNT(*) FROM questions WHERE quiz_id = $1`
)
