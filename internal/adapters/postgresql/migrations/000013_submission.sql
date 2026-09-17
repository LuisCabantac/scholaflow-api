-- +goose Up
CREATE TABLE IF NOT EXISTS public."submission" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    stream_id UUID NOT NULL REFERENCES public."stream"(id) ON DELETE CASCADE,
    classroom_id UUID NOT NULL REFERENCES public."classroom"(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    attachments TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
    links TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
    content TEXT,
    grade INTEGER,
    is_turned_in BOOLEAN NOT NULL DEFAULT false,
    is_graded BOOLEAN NOT NULL DEFAULT false,
    is_returned BOOLEAN NOT NULL DEFAULT false,
    submitted_at TIMESTAMP WITHOUT TIME ZONE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_stream_user_submission UNIQUE (stream_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_submission_stream_id ON public."submission"(stream_id);
CREATE INDEX IF NOT EXISTS idx_submission_user_classroom ON public."submission"(classroom_id, user_id);

-- +goose Down
DROP TABLE IF EXISTS public."submission";
