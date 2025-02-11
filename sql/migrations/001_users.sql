-- +goose Up

CREATE TABLE users (
    id uuid primary key ,
    username VARCHAR(100) NOT NULL unique ,
    email VARCHAR(255) UNIQUE NOT NULL,
    is_verified boolean default false,
    verify_code text null ,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    refresh_token text null
);

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE likes (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_id UUID NOT NULL,
    entity_type TEXT NOT NULL CHECK (entity_type IN ('post', 'comment')),
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT unique_like UNIQUE (user_id, entity_id, entity_type)
);

CREATE TABLE comments (
    id UUID PRIMARY KEY,
    entity_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    entity_type TEXT NOT NULL,
    is_edited boolean default false,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);


-- +goose Down

DROP TABLE comments;
DROP TABLE likes;
DROP TABLE posts;
DROP TABLE users;

