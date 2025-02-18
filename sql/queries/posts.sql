-- name: CreatePost :one
INSERT INTO posts (
    id, author_id, title, content, created_at, updated_at
) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;

-- name: UpdatePost :one
UPDATE posts SET title = $2, content = $3, updated_at = $4 WHERE id = $1 RETURNING *;



-- name: GetPost :one
SELECT
    p.id,
    p.author_id,
    p.title,
    p.content,
    p.created_at,
    p.updated_at,
    COALESCE(l.like_count, 0) AS like_count,
    COALESCE(lb.liked_by_user, false) AS liked_by_user,
    COALESCE(cc.comment_count, 0) AS comment_count
FROM posts p
         LEFT JOIN (
    SELECT
        post_id,
        COUNT(*) AS like_count
    FROM post_likes
    GROUP BY post_id
) AS l ON p.id = l.post_id

         LEFT JOIN (
    SELECT
        post_id,
        true AS liked_by_user
    FROM post_likes
    WHERE post_likes.user_id = $1
) AS lb ON p.id = lb.post_id

         LEFT JOIN (
    SELECT
        post_id,
        COUNT(*) AS comment_count
    FROM comments
    GROUP BY post_id
) AS cc ON p.id = cc.post_id

WHERE p.id = $1;



-- name: GetPosts :many
SELECT
    p.id,
    p.author_id,
    p.title,
    p.content,
    p.created_at,
    p.updated_at,
    COALESCE(l.like_count, 0) AS like_count,
    COALESCE(lb.liked_by_user, false) AS liked_by_user,
    COALESCE(cc.comment_count, 0) AS comment_count
FROM posts p
         LEFT JOIN (
    SELECT
        post_id,
        COUNT(*) AS like_count
    FROM post_likes
    GROUP BY post_id
) AS l ON p.id = l.post_id

         LEFT JOIN (
    SELECT
        post_id,
        true AS liked_by_user
    FROM post_likes
    WHERE post_likes.user_id = $1
) AS lb ON p.id = lb.post_id

         LEFT JOIN (
    SELECT
        post_id,
        COUNT(*) AS comment_count
    FROM comments
    GROUP BY post_id
) AS cc ON p.id = cc.post_id;

