CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS params (
    user_uuid UUID PRIMARY KEY,
    lang_code VARCHAR(2) NOT NULL CHECK (lang_code IN ('RU', 'EN')) DEFAULT 'RU',
    sound_volume INTEGER NOT NULL CHECK (sound_volume BETWEEN 0 AND 100) DEFAULT 100,
    game_sound_enabled BOOLEAN NOT NULL DEFAULT true
);