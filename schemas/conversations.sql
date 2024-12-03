CREATE TABLE IF NOT EXISTS conversations (
  id                  TEXT PRIMARY KEY,
  created_at          DATETIME,
  type                TEXT NOT NULL CHECK (type in ('direct', 'group')),
  user_id             TEXT,
  last_sent_user_id   TEXT,
  other_user_id       TEXT
)
