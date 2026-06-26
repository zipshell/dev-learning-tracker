-- name: FindUserById :one
SELECT *
FROM users
WHERE users.id = $1;

-- name: CreateUser :one
INSERT INTO users (email, password)
VALUES ($1, $2)
RETURNING *;

-- name: ListFoldersByUserId :many
SELECT * 
FROM folders INNER JOIN user_folders on folders.id = user_folders.folder_id
WHERE user_folders.user_id = $1
ORDER BY name;

-- name: FindFolderById :one
SELECT * FROM folders WHERE id = $1;

-- name: CreateFolder :one
INSERT INTO folders (name, description, parent_folder_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AssignFolderToUser :one
INSERT INTO user_folders (user_id, folder_id)
VALUES ($1, $2)
RETURNING *;

-- name: FindUsersIdsByFolderId :many
SELECT * FROM user_folders
WHERE user_folders.folder_id = $1;

-- name: ListEntriesByFolderId :many
SELECT * 
FROM entries
WHERE entries.folder_id = $1
ORDER BY name;

-- name: FindEntryById :one
SELECT * FROM entries WHERE entries.id = $1;

-- name: CreateEntry :one
INSERT INTO entries (name, content, folder_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateSession :one
INSERT INTO sessions (token, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: FindUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: FindSessionByToken :one
SELECT *
FROM sessions
WHERE token = $1;

-- name: FindActiveSessionByToken :one
SELECT *
FROM sessions
WHERE token = $1
  AND expired_at > NOW();

-- name: FindSessionsByUserId :many
SELECT *
FROM sessions
WHERE user_id = $1
ORDER BY expired_at ASC;

-- name: DeleteSessionsByIds :exec
DELETE FROM sessions
WHERE id = ANY($1::bigint[]);

-- name: DeleteSessionByToken :exec
DELETE FROM sessions
WHERE token = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expired_at < NOW();