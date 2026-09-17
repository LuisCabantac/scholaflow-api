-- +goose Up
CREATE TABLE IF NOT EXISTS public."classroom_message" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    classroom_id UUID NOT NULL REFERENCES public."classroom"(id) ON DELETE CASCADE,
    message TEXT,
    attachments TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE INDEX IF NOT EXISTS idx_classroom_message_feed ON public."classroom_message"(classroom_id, created_at);

-- +goose Down
DROP TABLE IF EXISTS public."classroom_message";
