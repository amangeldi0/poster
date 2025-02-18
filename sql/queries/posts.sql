-- name: GetPost :one
SELECT * FROM posts WHERE id = $1;

-- name: GetPosts :many
SELECT * FROM posts ORDER BY created_at DESC;

-- name: GetPostsWithLikes :many
SELECT
    p.id,
    p.author_id,
    p.title,
    p.content,
    p.created_at,
    p.updated_at,
    COALESCE(l.like_count, 0) AS like_count,
    COALESCE(lb.liked_by_user, false) AS liked_by_user
FROM posts p
         LEFT JOIN (
    SELECT entity_id, COUNT(*) AS like_count
    FROM likes
    WHERE entity_type = 'post'
    GROUP BY entity_id
) l ON p.id = l.entity_id
         LEFT JOIN (
    SELECT l.entity_id, true AS liked_by_user
    FROM likes l
    WHERE l.user_id = $1 AND l.entity_type = 'post'
) lb ON p.id = lb.entity_id
ORDER BY p.created_at DESC;


-- name: CreatePost :one
INSERT INTO posts (
    id, author_id, title, content, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;

-- name: UpdatePost :one
UPDATE posts SET title = $2, content = $3, updated_at = $4 WHERE id = $1 RETURNING *;





