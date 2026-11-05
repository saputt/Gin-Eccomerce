CREATE TABLE IF NOT EXISTS inventories (
    id VARCHAR(36) PRIMARY KEY,
    stock INT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    product_id VARCHAR(36),
    CONSTRAINT fk_inventory_product FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE CASCADE ON DELETE CASCADE
)