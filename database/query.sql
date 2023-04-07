-- name: GetMessage :one
SELECT * FROM messages WHERE id = ?;

-- name: ListMessages :many
SELECT * FROM messages WHERE conversation_id = ? ORDER BY id;

-- name: CreateMessage :one
INSERT INTO messages (conversation_id, content, author) VALUES (?, ?, ?) RETURNING *;

-- name: AppendMessage :one
UPDATE messages SET content = content || ? WHERE id = ? RETURNING *;

-- name: ListConversations :many
SELECT * FROM conversations ORDER BY last_message_time DESC;

-- name: CreateConversation :one
INSERT INTO conversations (title, last_message_time) VALUES (?, ?) RETURNING *;

-- name: DeleteConversation :exec
DELETE FROM conversations WHERE id = ?;
