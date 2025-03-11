CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.users (
                                user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                username VARCHAR(50) UNIQUE NOT NULL,
                                email VARCHAR(255) UNIQUE NOT NULL,
                                password_hash VARCHAR(255) NOT NULL,
                                role VARCHAR(10) CHECK (role IN ('buyer', 'artist', 'both')) NOT NULL,
                                created_at TIMESTAMP DEFAULT NOW(),
                                updated_at TIMESTAMP DEFAULT NOW(),
                                avatar_url VARCHAR(255),
                                bio TEXT
);

CREATE TABLE IF NOT EXISTS  users.user_reviews (
                                review_id SERIAL PRIMARY KEY,
                                reviewer_id UUID NOT NULL REFERENCES users.users(user_id) ON DELETE CASCADE,
                                reviewee_id UUID NOT NULL REFERENCES users.users(user_id) ON DELETE CASCADE,
                                rating INT CHECK (rating BETWEEN 1 AND 5) NOT NULL,
                                comment TEXT,
                                created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS  idx_users_email ON users.users(email);
CREATE INDEX IF NOT EXISTS  idx_users_username ON users.users(username);
CREATE INDEX IF NOT EXISTS  idx_reviews_reviewee ON users.user_reviews(reviewee_id);

CREATE SCHEMA IF NOT EXISTS drawings;

CREATE TABLE IF NOT EXISTS drawings.drawings (
                                                 drawing_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                                 artist_id UUID NOT NULL REFERENCES users.users(user_id) ON DELETE CASCADE,
                                                 file_path TEXT NOT NULL,
                                                 storage_provider TEXT NOT NULL CHECK (storage_provider IN ('local', 's3', 'other')),
                                                 visibility TEXT NOT NULL CHECK (visibility IN ('private', 'public', 'link')),
                                                 created_at TIMESTAMP DEFAULT NOW(),
                                                 updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS drawings.tags (
                                             tag_id SERIAL PRIMARY KEY,
                                             tag_name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS drawings.drawings_tags (
                                                      drawing_id UUID NOT NULL REFERENCES drawings.drawings(drawing_id) ON DELETE CASCADE,
                                                      tag_id INT NOT NULL REFERENCES drawings.tags(tag_id) ON DELETE CASCADE,
                                                      PRIMARY KEY (drawing_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_drawings_artist ON drawings.drawings(artist_id);
CREATE INDEX IF NOT EXISTS idx_tags_name ON drawings.tags(tag_name);
CREATE INDEX IF NOT EXISTS idx_drawings_tags ON drawings.drawings_tags(tag_id);
