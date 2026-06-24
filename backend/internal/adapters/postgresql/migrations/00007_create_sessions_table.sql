-- +goose Up
CREATE TABLE IF NOT EXISTS sessions (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL DEFAULT now() + interval '7 days',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);

-- +goose Down
DROP TABLE IF EXISTS sessions;
