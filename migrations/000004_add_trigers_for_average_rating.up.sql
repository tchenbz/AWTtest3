CREATE OR REPLACE FUNCTION update_average_rating() RETURNS TRIGGER AS $$
BEGIN
    UPDATE books
    SET average_rating = (SELECT AVG(rating) FROM reviews WHERE book_id = NEW.book_id)
    WHERE id = NEW.book_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_average_rating
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH ROW EXECUTE FUNCTION update_average_rating();
