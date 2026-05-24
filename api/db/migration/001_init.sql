CREATE TYPE book_status AS ENUM ('unread', 'reading', 'read', 'lent');

CREATE TABLE IF NOT EXISTS books (
    id             UUID        NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    isbn           TEXT        NOT NULL,
    title          TEXT        NOT NULL,
    author         TEXT,
    publisher      TEXT,
    published_date DATE,
    thumbnail_url  TEXT,
    status         book_status NOT NULL DEFAULT 'unread',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_books_created_at_id ON books (created_at DESC, id DESC);

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER books_updated_at
    BEFORE UPDATE ON books
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
