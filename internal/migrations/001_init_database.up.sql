CREATE TYPE role_type AS ENUM ('user', 'admin');
CREATE TYPE ad_status AS ENUM ('pending', 'approved', 'rejected');
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(150),
    last_name VARCHAR(150),
    phone VARCHAR(15) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role role_type NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens(
    user_id UUID REFERENCES users(id),
    token TEXT UNIQUE,
    expires_at TIMESTAMP NOT NULL
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ads (
    id SERIAL PRIMARY KEY,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location VARCHAR(100),
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    status ad_status NOT NULL DEFAULT 'pending',
    rejection_reason TEXT DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT false, -- send for moderation to make it true
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (
        (status = 'rejected' AND rejection_reason IS NOT NULL) OR
        (status IN ('pending', 'approved') AND rejection_reason IS NULL)
    )
);

CREATE TABLE ad_files (
    id SERIAL PRIMARY KEY,
    ad_id INTEGER NOT NULL REFERENCES ads(id) ON DELETE CASCADE,
    url VARCHAR(512) NOT NULL,
    file_name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);