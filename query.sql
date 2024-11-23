-- name: CreateUpload :one
INSERT INTO uploads (id) VALUES (?) RETURNING *;
