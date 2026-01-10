DROP TABLE IF EXISTS users CASCADE;

DROP TABLE IF EXISTS orders CASCADE;

DROP INDEX IF EXISTS idx_orders_new;
DROP INDEX IF EXISTS idx_orders_processing_ready;
DROP INDEX IF EXISTS idx_orders_user;

DROP TABLE IF EXISTS transactions CASCADE;

DROP INDEX IF EXISTS transactions_order_number_unique;
DROP INDEX IF EXISTS idx_transactions_user_withdrawals;
DROP INDEX IF EXISTS idx_transactions_user;