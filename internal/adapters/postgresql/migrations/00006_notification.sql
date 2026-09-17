-- +goose Up
-- -- +goose StatementBegin
DO $$ BEGIN
    CREATE TYPE public.notification_type AS ENUM ('stream', 'assignment', 'quiz', 'question', 'material', 'comment', 'join', 'addToClass', 'submit', 'grade');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS public."notification" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    actor_id TEXT REFERENCES public."user"(id) ON DELETE SET NULL,
    type public.notification_type NOT NULL,
    resource_id TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE INDEX IF NOT EXISTS idx_notification_user_id ON public."notification"(user_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS public."notification";
