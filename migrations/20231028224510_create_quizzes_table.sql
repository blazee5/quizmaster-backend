-- +goose Up
-- +goose StatementBegin
CREATE TABLE quizzes(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    user_id int NOT NULL,
    created_at DATE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE quizzes;
-- +goose StatementEnd
