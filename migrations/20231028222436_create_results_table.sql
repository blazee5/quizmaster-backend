-- +goose Up
-- +goose StatementBegin
CREATE TABLE results(
    id          SERIAL PRIMARY KEY,
    user_id     INT     NOT NULL,
    quiz_id     INT     NOT NULL,
    question_id INT     NOT NULL,
    answer_id   INT     NOT NULL,
    is_correct  BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id),
    FOREIGN KEY (question_id) REFERENCES questions (id),
    FOREIGN KEY (answer_id) REFERENCES answers (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE results;
-- +goose StatementEnd
