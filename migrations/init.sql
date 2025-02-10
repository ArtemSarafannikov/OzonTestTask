CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    last_activity TIMESTAMP DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_users_login ON users(login);

CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author_id UUID REFERENCES users(id) NOT NULL,
    title VARCHAR NOT NULL,
    content VARCHAR NOT NULL,
    allowed_comments BOOLEAN DEFAULT TRUE,
    edited_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_posts_author_id ON posts(author_id);

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID REFERENCES posts(id) NOT NULL,
    parent_comment_id UUID REFERENCES comments(id),
    author_id UUID REFERENCES users(id) NOT NULL,
    text VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);