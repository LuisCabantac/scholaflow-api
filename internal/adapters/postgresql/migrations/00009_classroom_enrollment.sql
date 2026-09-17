-- +goose Up
CREATE TABLE IF NOT EXISTS public."classroom_enrollment" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    classroom_id UUID NOT NULL REFERENCES public."classroom"(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_classroom_enrollment UNIQUE (classroom_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_classroom_enrollment_user_id ON public."classroom_enrollment"(user_id);
CREATE INDEX IF NOT EXISTS idx_classroom_enrollment_classroom_id ON public."classroom_enrollment"(classroom_id);

-- +goose Down
DROP TABLE IF EXISTS public."classroom_enrollment";
