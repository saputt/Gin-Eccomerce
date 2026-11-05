ALTER TABLE orders DROP CONSTRAINT fk_order_buyer;
ALTER TABLE orders DROP CONSTRAINT fk_order_seller;
ALTER TABLE orders DROP COLUMN seller_id;
ALTER TABLE orders DROP COLUMN buyer_id;
