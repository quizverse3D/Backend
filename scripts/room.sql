CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL,
    name VARCHAR(128) NOT NULL,
    password_hash TEXT,
    max_players INT NOT NULL CHECK (max_players > 0 and max_players < 32),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_public BOOLEAN NOT NULL DEFAULT true
);

-- участников будем хранить в Redis