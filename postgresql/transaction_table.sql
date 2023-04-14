-- DO NOT MODIFY

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    transaction_type VARCHAR(20) NOT NULL,
    customer_number INTEGER NOT NULL,
    transaction_amount BIGINT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_type_customer ON transactions (transaction_type, customer_number);

CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions (timestamp);