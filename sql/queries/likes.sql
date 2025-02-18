-- name: LikeComment :execrows
INSERT INTO comment_likes (id, user_id, comment_id, created_at)
VALUES ($1, $2, $3, $4);

-- name: UnlikeComment :execrows
DELETE FROM comment_likes
WHERE user_id = $1 AND comment_id = $2;

-- name: LikePost :execrows
INSERT INTO post_likes (id, user_id, post_id, created_at)
VALUES ($1, $2, $3, $4);

-- name: UnlikePost :execrows
DELETE FROM post_likes
WHERE user_id = $1 AND post_id = $2;