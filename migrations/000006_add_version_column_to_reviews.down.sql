-- ./migrations/000003_add_version_column_to_reviews.down.sql
ALTER TABLE reviews
DROP COLUMN IF EXISTS version;
