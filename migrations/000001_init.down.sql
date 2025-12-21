DROP TABLE IF EXISTS users CASCADE;

DROP INDEX IF EXISTS idx_users_login;

DROP TABLE IF EXISTS orders CASCADE;

DROP INDEX IF EXISTS idx_orders_user_uploaded_at;
DROP INDEX IF EXISTS idx_orders_status_incomplete;

DROP TABLE IF EXISTS transactions CASCADE;

DROP INDEX IF EXISTS transactions_order_id_unique ;
