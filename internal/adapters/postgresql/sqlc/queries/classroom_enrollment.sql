-- name: CreateClassroomEnrollment :one
INSERT INTO public."classroom_enrollment" (
    classroom_id, user_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetClassroomEnrollment :one
SELECT * FROM public."classroom_enrollment"
WHERE id = $1;

-- name: DeleteEnrollmentByClassroomAndUserID :execrows
DELETE FROM public."classroom_enrollment"
WHERE classroom_id = $1 AND user_id = $2;
