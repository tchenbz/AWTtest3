-- ./migrations/000002_add_version_column_to_books.down.sql
ALTER TABLE books
DROP COLUMN IF EXISTS version;
