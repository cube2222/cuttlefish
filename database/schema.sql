CREATE TABLE IF NOT EXISTS messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  conversation_id INTEGER NOT NULL,
  content TEXT NOT NULL,
  author TEXT NOT NULL, -- 'user', 'assistant', or tool name
  FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS conversations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  conversation_settings_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  last_message_time DATETIME NOT NULL,
  generating BOOLEAN NOT NULL DEFAULT 0,
  FOREIGN KEY (conversation_settings_id) REFERENCES conversation_settings(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS conversation_settings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  system_prompt_template TEXT NOT NULL DEFAULT 'You are a helpful assistant. Respond to the queries as best as you can.',
  tools_enabled TEXT_ARRAY NOT NULL DEFAULT '[]'
);

CREATE TABLE IF NOT EXISTS key_values (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS conversation_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    conversation_settings_id INTEGER NOT NULL,
    FOREIGN KEY (conversation_settings_id) REFERENCES conversation_settings(id) ON DELETE CASCADE ON UPDATE CASCADE
);
