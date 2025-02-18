-- name: LikeEntity :execrows
INSERT INTO likes (id, user_id, entity_id, entity_type, created_at)
SELECT $1, $2, $3, $4, $5
WHERE EXISTS (
    SELECT 1 FROM posts p WHERE p.id = $3 AND $4 = 'post'
    UNION ALL
    SELECT 1 FROM comments c WHERE c.id = $3 AND $4 = 'comment'
);

-- name: UnlikeEntity :execrows
DELETE FROM likes
WHERE likes.user_id = $1
  AND likes.entity_id = $2
  AND likes.entity_type = $3
  AND EXISTS (
    SELECT 1 FROM posts p WHERE p.id = $2 AND $3 = 'post'
    UNION ALL
    SELECT 1 FROM comments c WHERE c.id = $2 AND $3 = 'comment'
);

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
        entity_id,
        COUNT(*) AS like_count
    FROM likes
    WHERE entity_type = 'comment'
    GROUP BY entity_id
    ) l ON c.id = l.entity_id

LEFT JOIN (
    SELECT DISTINCT ON (l.entity_id)
    l.entity_id, true AS liked_by_user
    FROM likes l
    WHERE l.user_id = $1 AND l.entity_type = 'comment'
) lb ON c.id = lb.entity_id

WHERE c.post_id = $2
ORDER BY c.created_at DESC;

