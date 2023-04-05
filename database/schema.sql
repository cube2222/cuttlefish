CREATE TABLE IF NOT EXISTS messages (
  id INTEGER PRIMARY KEY,
  conversation_id INTEGER NOT NULL,
  content TEXT NOT NULL,
  sent_by_self BOOLEAN NOT NULL,
  FOREIGN KEY (conversation_id) REFERENCES conversations(id)
);

CREATE TABLE IF NOT EXISTS conversations (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  last_message_time DATETIME NOT NULL,
  system_prompt TEXT NOT NULL DEFAULT 'You are a helpful assistant. Respond to the queries as best as you can.'
);