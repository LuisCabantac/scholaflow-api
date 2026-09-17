-- +goose Up
DO $$ BEGIN
    CREATE TYPE public.role AS ENUM ('user', 'admin');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS public."user" (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    image TEXT NOT NULL,
    role public.role NOT NULL DEFAULT 'user'::public.role,
    school_name TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS public."user";
DROP TYPE IF EXISTS public.role;
