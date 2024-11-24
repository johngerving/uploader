CREATE TABLE IF NOT EXISTS uploads (
    id TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS parts (
    id INT NOT NULL,
    upload_id TEXT NOT NULL,
    data BLOB,
    PRIMARY KEY(upload_id, id),
    FOREIGN KEY(upload_id) REFERENCES uploads(id)
);