-- ./migrations/000003_add_version_column_to_reviews.up.sql
ALTER TABLE reviews
ADD COLUMN version INT DEFAULT 1;
