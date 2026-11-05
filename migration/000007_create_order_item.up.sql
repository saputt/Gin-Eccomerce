CREATE TABLE IF NOT EXISTS order_items (
	id        VARCHAR(36) PRIMARY KEY,
	price  FLOAT,
	quantity  INT,
	total     FLOAT,

	product_id VARCHAR(36),
	order_id   VARCHAR(36),

    CONSTRAINT fk_orderitem_product FOREIGN KEY (product_id) REFERENCES products(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_orderitem_order FOREIGN KEY (order_id) REFERENCES orders(id) ON UPDATE CASCADE ON DELETE CASCADE
)