CREATE TABLE IF NOT EXISTS credentials (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    hash_algorithm TEXT NOT NULL
);