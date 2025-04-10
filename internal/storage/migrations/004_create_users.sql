CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('employee', 'moderator')),
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);