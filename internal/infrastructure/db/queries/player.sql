-- name: FindAllPlayers :many
SELECT * FROM player
WHERE active = 1
ORDER BY elo_rating DESC;

-- name: FindPlayerByID :one
SELECT * FROM player
WHERE id = ? LIMIT 1;

-- name: FindPlayerByNickname :one
SELECT * FROM player
WHERE nickname = ? COLLATE NOCASE LIMIT 1;

-- name: CreatePlayer :one
INSERT INTO player (id, nickname, elo_rating)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdatePlayer :one
UPDATE player
SET nickname = ?, elo_rating = ?, active = ?
WHERE id = ?
RETURNING *;

-- name: UntrackPlayer :one
UPDATE player
SET active = 0
WHERE id = ?
RETURNING *;

-- name: DeletePlayer :exec
DELETE FROM player
WHERE id = ?;