DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_product') THEN
        CREATE TYPE status_product AS ENUM ('ACTIVE', 'ARCHIVE');
    END IF;
END $$;

ALTER TABLE products ADD status status_product DEFAULT 'ACTIVE'