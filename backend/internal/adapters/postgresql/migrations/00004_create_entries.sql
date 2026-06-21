-- +goose Up
CREATE TABLE IF NOT EXISTS entries (
	id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    content TEXT,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	folder_id BIGINT NOT NULL REFERENCES folders(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_entries_folder_id ON entries(folder_id);

-- +goose Down
DROP TABLE IF EXISTS entries;
