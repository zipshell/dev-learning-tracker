-- +goose Up
CREATE TABLE IF NOT EXISTS folders (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	parent_folder_id BIGINT REFERENCES folders(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_folders_parent_folder_id ON folders(parent_folder_id);

-- +goose Down
DROP TABLE IF EXISTS folders;
