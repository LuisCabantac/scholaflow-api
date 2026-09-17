-- +goose Up
CREATE TABLE IF NOT EXISTS public."note" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    title TEXT,
    content TEXT,
    attachments TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
    is_pinned BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE INDEX IF NOT EXISTS idx_note_user_pinned ON public."note"(user_id, is_pinned DESC, updated_at DESC);

-- +goose Down
DROP TABLE IF EXISTS public."note";
