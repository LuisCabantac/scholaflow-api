-- +goose Up
CREATE TABLE IF NOT EXISTS public."classroom_topic" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    classroom_id UUID NOT NULL REFERENCES public."classroom"(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    CONSTRAINT uq_classroom_topic_name UNIQUE (classroom_id, name)
);

CREATE INDEX IF NOT EXISTS idx_classroom_topic_classroom_id ON public."classroom_topic"(classroom_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS public."classroom_topic";
