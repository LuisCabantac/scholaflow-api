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
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_account_user_id ON public."account"(user_id);

-- +goose Down
DROP TABLE IF EXISTS public."account";
