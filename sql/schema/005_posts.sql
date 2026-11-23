-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title VARCHAR(255),
    url VARCHAR(255) UNIQUE NOT NULL,
    description VARCHAR(5000) NULL,
    published_at TIMESTAMP NULL,
    feed_id UUID REFERENCES feeds (id) NOT NULL
);

-- +goose Down
DROP TABLE posts;
