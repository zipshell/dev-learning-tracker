-- +goose Up
CREATE TABLE IF NOT EXISTS entry_tags (
	id BIGSERIAL PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	entry_id BIGINT NOT NULL REFERENCES entries(id) ON DELETE CASCADE,
	tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_entry_tags_entry_tag ON entry_tags(entry_id, tag_id);
CREATE INDEX IF NOT EXISTS idx_entry_tags_tag_id ON entry_tags(tag_id);

-- +goose Down
DROP TABLE IF EXISTS entry_tags;
