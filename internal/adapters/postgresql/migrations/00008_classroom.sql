-- +goose Up
CREATE TABLE IF NOT EXISTS public."classroom" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    subject TEXT,
    section TEXT NOT NULL,
    description TEXT,
    room TEXT,
    code TEXT NOT NULL UNIQUE,
    card_background TEXT NOT NULL DEFAULT '#a7adcb',
    illustration_index INTEGER NOT NULL DEFAULT 0,
    allow_users_to_comment BOOLEAN NOT NULL DEFAULT false,
    allow_users_to_post BOOLEAN NOT NULL DEFAULT false,
    teacher_id TEXT NOT NULL REFERENCES public."user"(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

CREATE INDEX IF NOT EXISTS idx_class_teacher_id ON public."classroom"(teacher_id);

-- +goose Down
DROP TABLE IF EXISTS public."classroom";
