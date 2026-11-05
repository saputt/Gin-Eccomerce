ALTER TABLE order_items DROP CONSTRAINT fk_orderitem_order;
ALTER TABLE orders DROP COLUMN status;
DROP TYPE IF EXISTS order_status;
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'suborder_status') THEN
        CREATE TYPE suborder_status AS ENUM ('CANCEL', 'PROCESS','SHIP','COMPLETED', 'PENDING');
    END IF;
END $$;
ALTER TABLE order_items DROP COLUMN order_id;
CREATE TABLE suborders (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36),
    seller_id VARCHAR(36) NOT NULL,
    status suborder_status DEFAULT 'PENDING',
    subtotal DECIMAL(12, 2) NOT NULL,
    shipping_cost DECIMAL(12, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_suborder_order FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE CASCADE ON DELETE CASCADE
);
ALTER TABLE order_items ADD COLUMN suborder_id VARCHAR(36);
ALTER TABLE order_items ADD CONSTRAINT fk_order_item_suborders FOREIGN KEY (suborder_id) REFERENCES suborders(id) ON UPDATE CASCADE ON DELETE CASCADE;