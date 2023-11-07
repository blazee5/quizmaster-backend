-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_answers(
    id          SERIAL PRIMARY KEY,
    user_id     INT     NOT NULL,
    question_id INT     NOT NULL,
    answer_id   INT     NOT NULL,
    text VARCHAR(255) NOT NULL DEFAULT '',
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (question_id) REFERENCES questions (id),
    FOREIGN KEY (answer_id) REFERENCES answers (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_answers;
-- +goose StatementEnd
