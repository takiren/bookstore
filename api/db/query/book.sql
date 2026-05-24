-- name: ListBooksFirst :many
SELECT id, isbn, title, author, publisher, published_date, thumbnail_url, status, created_at, updated_at
FROM books
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(limit_count)::INT;

-- name: ListBooksAfterCursor :many
SELECT id, isbn, title, author, publisher, published_date, thumbnail_url, status, created_at, updated_at
FROM books
WHERE (created_at < sqlc.arg(cursor_created_at)::TIMESTAMPTZ)
   OR (created_at = sqlc.arg(cursor_created_at)::TIMESTAMPTZ AND id < sqlc.arg(cursor_id)::UUID)
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(limit_count)::INT;
