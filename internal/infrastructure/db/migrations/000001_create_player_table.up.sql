CREATE TABLE player (
    id TEXT PRIMARY KEY,
    nickname TEXT NOT NULL,
    elo_rating INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')), -- ISO 8601 format
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now')), -- ISO 8601 format
    active INTEGER NOT NULL DEFAULT 1 -- 1 for true, 0 for false
);

CREATE UNIQUE INDEX idx_player_nickname ON player(nickname);

CREATE TRIGGER update_player_updated_at
AFTER UPDATE ON player
BEGIN
    UPDATE player
    SET updated_at = strftime('%Y-%m-%dT%H:%M:%SZ', 'now')
    WHERE id = NEW.id;
END;

INSERT INTO player (id, nickname, elo_rating) VALUES
('1', 'parpi', 2200),
('2', 'fbr', 2400),
('3', 'milteira', 1800);
