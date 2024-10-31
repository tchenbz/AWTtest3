-- ./migrations/000002_add_version_column_to_products.up.sql
ALTER TABLE products
ADD COLUMN version INT DEFAULT 1;
