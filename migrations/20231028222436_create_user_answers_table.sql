-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_answers(
    id          SERIAL PRIMARY KEY,
    user_id     INT          NOT NULL,
    question_id INT          NOT NULL,
    answer_id   INT          NOT NULL,
    result_id   INT          NOT NULL,
    text        VARCHAR(255) NOT NULL DEFAULT '',
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions (id) ON DELETE CASCADE,
    FOREIGN KEY (answer_id) REFERENCES answers (id) ON DELETE CASCADE,
    FOREIGN KEY (result_id) REFERENCES results (id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_answers;
-- +goose StatementEnd
