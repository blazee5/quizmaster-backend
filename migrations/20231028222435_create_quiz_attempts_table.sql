-- +goose Up
-- +goose StatementBegin
CREATE TABLE quiz_attempts(
    id SERIAL PRIMARY KEY,
    quiz_id INT NOT NULL,
    user_id int NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    FOREIGN KEY (quiz_id) REFERENCES quizzes (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE quiz_attempts;
-- +goose StatementEnd
