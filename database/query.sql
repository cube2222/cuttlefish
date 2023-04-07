-- name: GetMessage :one
SELECT * FROM messages WHERE id = ?;

-- name: ListMessages :many
SELECT * FROM messages WHERE conversation_id = ? ORDER BY id;

-- name: CreateMessage :one
INSERT INTO messages (conversation_id, content, author) VALUES (?, ?, ?) RETURNING *;

-- name: AppendMessage :one
UPDATE messages SET content = content || ? WHERE id = ? RETURNING *;

-- name: GetConversation :one
SELECT * FROM conversations WHERE id = ?;

-- name: ListConversations :many
SELECT * FROM conversations ORDER BY last_message_time DESC;

-- name: CreateConversation :one
INSERT INTO conversations (conversation_settings_id, title, last_message_time) VALUES (?, ?, ?) RETURNING *;

-- name: DeleteConversation :exec
DELETE FROM conversations WHERE id = ?;

-- name: MarkGenerationStarted :exec
UPDATE conversations SET generating = true WHERE id = ?;

-- name: MarkGenerationDone :exec
UPDATE conversations SET generating = false WHERE id = ?;

-- name: GetConversationSettings :one
SELECT * FROM conversation_settings WHERE id = ?;

-- name: CreateConversationSettings :one
INSERT INTO conversation_settings (system_prompt_template, tools_enabled) VALUES (?, ?) RETURNING *;

-- name: UpdateConversationSettings :exec
UPDATE conversation_settings SET system_prompt_template = ?, tools_enabled = ? WHERE id = ?;

-- name: GetKeyValue :one
SELECT * FROM key_values WHERE key = ?;

-- name: SetKeyValue :exec
INSERT INTO key_values (key, value) VALUES (?, ?) ON CONFLICT (key) DO UPDATE SET value = ?;
-- TODO: Not sure if this will actually work.
