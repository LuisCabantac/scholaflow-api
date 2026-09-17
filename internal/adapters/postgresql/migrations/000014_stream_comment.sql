-- +goose Up
CREATE TABLE IF NOT EXISTS public."stream_comment" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES public."stream"(id) ON DELETE CASCADE,
    classroom_id UUID NOT NULL REFERENCES public."classroom"(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    content TEXT,
    attachment TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE INDEX IF NOT EXISTS idx_stream_comment_feed ON public."stream_comment"(stream_id, created_at ASC);

-- +goose Down
DROP TABLE IF EXISTS public."stream_comment";
