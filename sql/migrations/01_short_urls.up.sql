CREATE TABLE IF NOT EXISTS short_urls (
    id           SERIAL PRIMARY KEY,
    user_id      TEXT,
    slug         TEXT UNIQUE NOT NULL,
    original_url TEXT UNIQUE NOT NULL,
    deleted      BOOLEAN NOT NULL DEFAULT FALSE
);
