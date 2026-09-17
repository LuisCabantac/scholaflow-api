-- +goose Up
CREATE TABLE IF NOT EXISTS public."passkey" (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT,
    public_key TEXT NOT NULL,
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    credential_id TEXT NOT NULL,
    counter INTEGER NOT NULL,
    device_type TEXT NOT NULL,
    backed_up BOOLEAN NOT NULL,
    transports TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE,
    aaguid TEXT
);

-- +goose Down
DROP TABLE IF EXISTS public."passkey";
