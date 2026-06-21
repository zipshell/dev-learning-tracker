-- +goose Up
CREATE TABLE IF NOT EXISTS user_folders (
	id BIGSERIAL PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	folder_id BIGINT NOT NULL REFERENCES folders(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_user_folders_user_folder ON user_folders(user_id, folder_id);
CREATE INDEX IF NOT EXISTS idx_user_folders_folder_id ON user_folders(folder_id);

-- +goose Down
DROP TABLE IF EXISTS user_folders;
