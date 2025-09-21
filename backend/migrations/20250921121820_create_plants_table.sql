-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS plants (
    id          SERIAL PRIMARY KEY,
    author      VARCHAR(255) NOT NULL,
    image_data  TEXT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS plants;
-- +goose StatementEnd
