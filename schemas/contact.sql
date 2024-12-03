CREATE TABLE IF NOT EXISTS contact (
  id          TEXT PRIMARY KEY,
  username    TEXT,
  status      TEXT NOT NULL CHECK (status in ('active', 'inactive', 'away'))
)
