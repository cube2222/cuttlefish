-- name: ListMessages :many
SELECT * FROM messages WHERE conversation_id = ? ORDER BY id;

-- name: CreateMessage :one
INSERT INTO messages (conversation_id, content, sent_by_self) VALUES (?, ?, ?) RETURNING *;

-- name: AppendMessage :one
UPDATE messages SET content = content || ? WHERE id = ? RETURNING *;

-- name: DeleteMessages :exec
DELETE FROM messages;

-- name: ListConversations :many
SELECT * FROM conversations ORDER BY last_message_time DESC;

-- name: CreateConversation :one
INSERT INTO conversations (title, last_message_time) VALUES (?, ?) RETURNING *;
