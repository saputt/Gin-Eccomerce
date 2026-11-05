ALTER TABLE order_items DROP CONSTRAINT fk_order_item_suborders;
ALTER TABLE order_items DROP COLUMN suborder_id;
DROP TABLE IF EXISTS suborders;
ALTER TABLE order_items ADD COLUMN IF NOT EXISTS order_id VARCHAR(36);
ALTER TABLE order_items ADD CONSTRAINT fk_orderitem_order FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE orders ADD COLUMN status order_status;