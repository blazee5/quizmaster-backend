-- +goose Up
-- +goose StatementBegin
CREATE TABLE questions(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    image VARCHAR(255),
    quiz_id int NOT NULL,
    CONSTRAINT quiz_id FOREIGN KEY (quiz_id) REFERENCES quizzes (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE questions;
-- +goose StatementEnd
