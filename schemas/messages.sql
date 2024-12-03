CREATE TABLE IF NOT EXISTS messages (
  id                TEXT PRIMARY KEY,
  conversation_id   TEXT,
  sender_id         TEXT,
  content           TEXT,
  created_at        DATETIME DEFAULT CURRENT_TIMESTAMP
)
