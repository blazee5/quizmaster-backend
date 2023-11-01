-- +goose Up
-- +goose StatementBegin
CREATE TABLE answers(
    id SERIAL PRIMARY KEY,
    text VARCHAR(255) NOT NULL,
    question_id int NOT NULL,
    is_correct BOOLEAN NOT NULL,
    CONSTRAINT question_id FOREIGN KEY (question_id) REFERENCES questions (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE answers;
-- +goose StatementEnd