-- ./migrations/000002_add_version_column_to_books.up.sql
ALTER TABLE books
ADD COLUMN version INT DEFAULT 1;
