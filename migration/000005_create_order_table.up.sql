DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
        CREATE TYPE order_status AS ENUM ('CANCEL', 'PROCESS','SHIP','COMPLETED');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS orders (
	id       VARCHAR(36) PRIMARY KEY,
	status   order_status DEFAULT 'PROCESS',
	subtotal FLOAT,
	seller_id VARCHAR(36),
	buyer_id   VARCHAR(36),

	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_order_buyer FOREIGN KEY (buyer_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_order_seller FOREIGN KEY (seller_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
)
