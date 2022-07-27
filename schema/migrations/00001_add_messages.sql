-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA core;
CREATE TABLE core.messages
(
    uid        uuid PRIMARY KEY UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    created_at timestamp               NOT NULL DEFAULT NOW(),
    updated_at timestamp               NOT NULL DEFAULT NOW(),
    deleted_at timestamp,

    content    text                    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE core.messages;
DROP SCHEMA core;
-- +goose StatementEnd
