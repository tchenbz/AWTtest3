CREATE OR REPLACE FUNCTION update_average_rating() RETURNS TRIGGER AS $$
BEGIN
    UPDATE products
    SET average_rating = (SELECT AVG(rating) FROM reviews WHERE product_id = NEW.product_id)
    WHERE id = NEW.product_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_average_rating
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH ROW EXECUTE FUNCTION update_average_rating();
