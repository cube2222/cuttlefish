CREATE TABLE IF NOT EXISTS messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  conversation_id INTEGER NOT NULL,
  content TEXT NOT NULL,
  author TEXT NOT NULL, -- 'user', 'assistant', or tool name
  FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS conversations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  last_message_time DATETIME NOT NULL,
  generating BOOLEAN NOT NULL DEFAULT 0,
  system_prompt TEXT NOT NULL DEFAULT 'You are a helpful assistant. Respond to the queries as best as you can.'
);
