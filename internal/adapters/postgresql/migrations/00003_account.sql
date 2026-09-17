-- +goose Up
CREATE TABLE IF NOT EXISTS public."account" (
    id TEXT PRIMARY KEY NOT NULL,
    account_id TEXT NOT NULL,
    provider_id TEXT NOT NULL,
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    access_token TEXT,
    refresh_token TEXT,
    id_token TEXT,
    access_token_expires_at TIMESTAMP WITHOUT TIME ZONE,
    refresh_token_expires_at TIMESTAMP WITHOUT TIME ZONE,
    scope TEXT,
    password TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS public."account";
