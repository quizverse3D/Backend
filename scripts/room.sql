CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    name VARCHAR(128) NOT NULL,
    password_hash TEXT,
    max_players INT NOT NULL CHECK (max_players > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- участников будем хранить в Redis