package models

type Result struct {
	Id         int  `json:"id" db:"id"`
	UserId     int  `json:"user_id" db:"user_id"`
	QuizId     int  `json:"quiz_id" db:"quiz_id"`
	QuestionId int  `json:"question_id" db:"question_id"`
	AnswerId   int  `json:"answer_id" db:"answer_id"`
	IsCorrect  bool `json:"is_correct" db:"is_correct"`
}
