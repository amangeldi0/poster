-- name: GetPost :one
SELECT * FROM posts WHERE id = $1;

-- name: GetPosts :many
SELECT * FROM posts ORDER BY created_at DESC;

-- name: CreatePost :one
INSERT INTO posts (
    id, author_id, title, content, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;

-- name: UpdatePost :one
UPDATE posts SET title = $2, content = $3, updated_at = $4 WHERE id = $1 RETURNING *;





