CREATE TABLE IF NOT EXISTS chat_messages (
  id INTEGER PRIMARY KEY,
  content TEXT NOT NULL,
  sent_by_self BOOLEAN NOT NULL
);
