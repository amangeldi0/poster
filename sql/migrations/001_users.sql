-- +goose Up

CREATE TABLE users (
    id uuid primary key ,
    username VARCHAR(100) NOT NULL unique ,
    email VARCHAR(255) UNIQUE NOT NULL,
    is_verified boolean default false,
    verify_code text null ,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    refresh_token text null
);

CREATE TABLE likes (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    entity_id UUID NOT NULL,
    entity_type TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT unique_like UNIQUE (user_id, entity_id, entity_type)
);

CREATE TABLE comments (
    id UUID PRIMARY KEY ,
    entity_id UUID NOT NULL,
    entity_type TEXT NOT NULL,
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE posts (
    id uuid primary key,
    author_id uuid not null,
    title VARCHAR(100) NOT NULL,
    content text,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP DEFAULT NOW()
);
