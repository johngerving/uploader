-- name: CreateUpload :one
INSERT INTO uploads (id) VALUES (?) RETURNING *;

-- name: FindUploadById :one
SELECT * FROM uploads WHERE id=?;
