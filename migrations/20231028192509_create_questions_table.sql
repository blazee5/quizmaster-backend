-- +goose Up
-- +goose StatementBegin
CREATE TABLE questions(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) DEFAULT '',
    image VARCHAR(255) DEFAULT '',
    quiz_id int NOT NULL,
    type VARCHAR(255) NOT NULL DEFAULT 'choice',
    order_id INT NOT NULL,
    CONSTRAINT quiz_id FOREIGN KEY (quiz_id) REFERENCES quizzes (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE questions;
-- +goose StatementEnd
