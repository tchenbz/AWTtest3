-- ./migrations/000002_add_version_column_to_products.down.sql
ALTER TABLE products
DROP COLUMN IF EXISTS version;
