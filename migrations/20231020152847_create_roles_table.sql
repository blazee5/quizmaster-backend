-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
INSERT INTO roles (name) VALUES ('user'), ('admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE roles;
-- +goose StatementEnd
