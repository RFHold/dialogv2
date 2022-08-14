-- +goose Up
-- +goose StatementBegin
CREATE TABLE core.users
(
    uid        uuid PRIMARY KEY UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    created_at timestamp               NOT NULL DEFAULT NOW(),
    updated_at timestamp               NOT NULL DEFAULT NOW(),
    deleted_at timestamp,

    email      text,
    phone      text                    NOT NULL,
    full_name  text                    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE core.users;
-- +goose StatementEnd
