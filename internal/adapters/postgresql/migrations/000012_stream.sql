-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    CREATE TYPE public.stream_type AS ENUM ('stream', 'assignment', 'quiz', 'question', 'material');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS public."stream" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    classroom_id UUID NOT NULL REFERENCES public."classroom"(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    classroom_topic_id UUID REFERENCES public."classroom_topic"(id) ON DELETE SET NULL,
    type public.stream_type NOT NULL DEFAULT 'stream'::public.stream_type,
    title TEXT,
    content TEXT,
    attachments TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
    links TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
    points INTEGER,
    is_pinned BOOLEAN NOT NULL DEFAULT false,
    accepting_submissions BOOLEAN NOT NULL DEFAULT true,
    close_submissions_after_due_date BOOLEAN NOT NULL DEFAULT true,
    announce_to TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
    announce_to_all BOOLEAN NOT NULL DEFAULT true,
    due_at TIMESTAMP WITHOUT TIME ZONE,
    scheduled_at TIMESTAMP WITHOUT TIME ZONE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE INDEX IF NOT EXISTS idx_stream_classroom_feed ON public."stream"(classroom_id, is_pinned DESC, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_stream_topic_id ON public."stream"(classroom_topic_id);

-- +goose Down
DROP TABLE IF EXISTS public."stream";
DROP TYPE IF EXISTS public.stream_type;
