DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status') THEN
        CREATE TYPE payment_status AS ENUM ('PENDING','PAID','EXPIRED','FAILED');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36),
    
    xendit_invoice_id VARCHAR(255) UNIQUE NOT NULL, 
    xendit_payment_method VARCHAR(50), 
    xendit_payment_channel VARCHAR(50),
    
    amount DECIMAL(12, 2) NOT NULL,
    payment_url TEXT,
    
    status payment_status NOT NULL DEFAULT 'PENDING', 
    raw_response JSON, 
    expired_at TIMESTAMP,
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_xendit_invoice_id ON payments(xendit_invoice_id);