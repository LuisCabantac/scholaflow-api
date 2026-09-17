-- +goose Up
CREATE TABLE IF NOT EXISTS public."session" (
    id TEXT PRIMARY KEY NOT NULL,
    token TEXT NOT NULL UNIQUE,
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    ip_address TEXT,
    user_agent TEXT,

);

-- +goose Down
DROP TABLE IF EXISTS public."session";
