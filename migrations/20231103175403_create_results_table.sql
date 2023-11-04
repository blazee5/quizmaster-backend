-- +goose Up
-- +goose StatementBegin
CREATE TABLE results(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    quiz_id INT NOT NULL,
    score INT NOT NULL,
    percent INT NOT NULL,
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE results;
-- +goose StatementEnd
