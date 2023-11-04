-- +goose Up
-- +goose StatementBegin
CREATE TABLE results(
    id SERIAL PRIMARY KEY,
    attempt_id int NOT NULL,
    score INT NOT NULL,
    percent DECIMAL NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE results;
-- +goose StatementEnd
