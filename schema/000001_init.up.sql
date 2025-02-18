CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.accounts (
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

CREATE TABLE users.user_reviews (
                                review_id SERIAL PRIMARY KEY,
                                reviewer_id UUID NOT NULL REFERENCES users.accounts(user_id) ON DELETE CASCADE,
                                reviewee_id UUID NOT NULL REFERENCES users.accounts(user_id) ON DELETE CASCADE,
                                rating INT CHECK (rating BETWEEN 1 AND 5) NOT NULL,
                                comment TEXT,
                                created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users.accounts(email);
CREATE INDEX idx_users_username ON users.accounts(username);
CREATE INDEX idx_reviews_reviewee ON users.user_reviews(reviewee_id);