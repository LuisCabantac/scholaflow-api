-- +goose Up
CREATE TABLE IF NOT EXISTS public."submission_comment" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES public."submission"(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    content TEXT,
    attachment TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE INDEX IF NOT EXISTS idx_submission_comment_thread ON public."submission_comment"(submission_id, created_at ASC);

-- +goose Down
DROP TABLE IF EXISTS public."submission_comment";
