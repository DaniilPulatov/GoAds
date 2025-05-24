CREATE TYPE role_type AS ENUM ('user', 'admin');
CREATE TYPE ad_status AS ENUM ('draft', 'approved', 'rejected');

CREATE table users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(150),
    last_name VARCHAR(150),
    phone VARCHAR(15) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role role_type NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE table refresh_tokens(
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
    status ad_status NOT NULL DEFAULT 'draft',
    rejection_reason TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (
        (status = 'rejected' AND rejection_reason IS NOT NULL) OR
        (status IN ('draft', 'approved') AND rejection_reason IS NULL)
    )
);

CREATE TABLE ad_files (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    url VARCHAR(512) NOT NULL,
    ad_id INTEGER NOT NULL REFERENCES ads(id) ON DELETE CASCADE,
    file_name VARCHAR(255),
);