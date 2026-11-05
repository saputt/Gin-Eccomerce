ALTER TABLE orders ADD COLUMN seller_id ;
ALTER TABLE orders ADD COLUMN buyer_id;
ALTER TABLE orders ADD CONSTRAINT fk_order_buyer FOREIGN KEY (buyer_id) REFERENCES users(id);
ALTER TABLE orders ADD CONSTRAINT fk_order_seller FOREIGN KEY (seller_id) REFERENCES users(id);