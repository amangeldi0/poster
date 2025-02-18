-- name: CreateComment :one
INSERT INTO comments (id, post_id, user_id, is_edited, content, created_at, updated_at)
SELECT $1, $2, $3, $4, $5, $6, $7 WHERE EXISTS(SELECT 1 FROM posts WHERE posts.id = $2) RETURNING *;

-- name: DeleteComment :execrows
DELETE FROM comments WHERE id = $1 AND post_id = $2 AND user_id = $3;

-- name: GetCommentsForPost :many
SELECT
    c.id,
    c.post_id,
    c.user_id,
    c.content,
    c.created_at,
    c.updated_at,
    COALESCE(l.like_count, 0) AS like_count,
    COALESCE(lb.liked_by_user, false) AS liked_by_user
FROM comments c
LEFT JOIN (
SELECT
    comment_id,
    COUNT(*) AS like_count
FROM comment_likes
GROUP BY comment_id
) l ON c.id = l.comment_id

LEFT JOIN (
SELECT
    comment_id,
    true AS liked_by_user
FROM comment_likes
WHERE c.user_id = $1
) lb ON c.id = lb.comment_id

WHERE c.post_id = $2
ORDER BY c.created_at DESC;
