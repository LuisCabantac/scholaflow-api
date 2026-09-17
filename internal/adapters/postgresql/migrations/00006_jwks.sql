-- +goose Up
CREATE TABLE IF NOT EXISTS public."jwks" (
    id TEXT PRIMARY KEY NOT NULL,
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL,
    alg TEXT,
    crv TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITHOUT TIME ZONE
);

-- +goose Down
DROP TABLE IF EXISTS public."jwks";
