-- name: CreateUpload :one
INSERT INTO uploads (id, complete) VALUES (?, 0) RETURNING *;

-- name: FindUploadById :one
SELECT * FROM uploads WHERE id=?;

-- name: CreatePart :exec
INSERT INTO parts (upload_id, id, data) VALUES (?, ?, ?);

-- name: FindUploadPartsById :many
SELECT id from parts WHERE upload_id=? ORDER BY id ASC;