-- name: ListMessages :many
SELECT * FROM chat_messages ORDER BY id;

-- name: CreateMessage :one
INSERT INTO chat_messages (content, sent_by_self) VALUES (?, ?) RETURNING *;

-- name: AppendMessage :one
UPDATE chat_messages SET content = content || ? WHERE id = ? RETURNING *;

-- name: DeleteMessages :exec
DELETE FROM chat_messages;
