-- +goose Up
-- +goose StatementBegin
CREATE TABLE verification_codes (
    id SERIAL PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    code VARCHAR(255) UNIQUE NOT NULL ,
    user_id INT NOT NULL,
    email VARCHAR(255) NOT NULL,
    expire_date TIMESTAMP DEFAULT NOW() + INTERVAL '10 hours',
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE verification_codes;
-- +goose StatementEnd
