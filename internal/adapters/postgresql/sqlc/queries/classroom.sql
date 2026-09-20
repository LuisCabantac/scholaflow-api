-- name: CreateClassroom :one
INSERT INTO public."classroom" (
    name, subject, section, description, code, room, card_background, illustration_index, teacher_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetClassroom :one
SELECT * FROM public."classroom"
WHERE id = $1;

-- name: GetClassroomByCode :one
SELECT * FROM public."classroom"
WHERE code = $1;
