-- name: CreateFeedback :one
INSERT INTO feedback (message)
VALUES ($1)
RETURNING id, message, created_at;

-- name: GetFeedback :one
SELECT id, message, created_at
FROM feedback
WHERE id = $1;

-- name: ListFeedback :many
SELECT id, message, created_at
FROM feedback
ORDER BY id DESC;

-- name: UpdateFeedback :one
UPDATE feedback
SET message = $2
WHERE id = $1
RETURNING id, message, created_at;

-- name: DeleteFeedback :exec
DELETE FROM feedback
WHERE id = $1;