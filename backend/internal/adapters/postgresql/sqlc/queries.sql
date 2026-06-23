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

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens
SET token = $1,
    expiration = $2
WHERE id = $3
RETURNING *;

-- name: FindUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: FindRefreshTokenByToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1;

-- name: FindRefreshTokensByUserId :many
SELECT *
FROM refresh_tokens
WHERE user_id = $1
ORDER BY expiration ASC;

-- name: DeleteRefreshTokensByIds :exec
DELETE FROM refresh_tokens
WHERE id = ANY($1::bigint[]);