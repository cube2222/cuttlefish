-- TODO: Change all wildcards to explicit column lists.

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

-- name: GetDefaultConversationSettings :one
SELECT * FROM conversation_settings WHERE is_default = true;

-- name: CreateConversationSettings :one
INSERT INTO conversation_settings (system_prompt_template, tools_enabled) VALUES (?, ?) RETURNING *;

-- name: UpdateConversationSettings :one
UPDATE conversation_settings SET system_prompt_template = ?, tools_enabled = ? WHERE id = ? RETURNING *;

-- name: CreateDefaultConversationSettings :one
INSERT INTO conversation_settings (system_prompt_template, tools_enabled, is_default) VALUES (?, ?, true) RETURNING *;

-- Doesn't work...
-- -- name: SetDefaultConversationSettings :exec
-- INSERT INTO conversation_settings (system_prompt_template, tools_enabled, is_default) VALUES (@systemprompttemplate, @toolsenabled, true) ON CONFLICT (is_default) DO UPDATE SET system_prompt_template = @systemprompttemplate, tools_enabled = @toolsenabled;

-- name: GetKeyValue :one
SELECT * FROM key_values WHERE key = ?;

-- name: CreateKeyValue :exec
INSERT INTO key_values (key, value) VALUES (?, ?);

-- name: UpdateKeyValue :exec
UPDATE key_values SET value = ? WHERE key = ?;

-- name: CloneConversationSettings :one
INSERT INTO conversation_settings(system_prompt_template, tools_enabled) SELECT system_prompt_template, tools_enabled FROM conversation_settings WHERE conversation_settings.id = ? RETURNING *;

-- name: CreateConversationTemplate :one
INSERT INTO conversation_templates(name, conversation_settings_id) VALUES (?, ?) RETURNING *;

-- name: ResetConversationFrom :exec
DELETE FROM messages WHERE conversation_id = ? AND id > ?;
