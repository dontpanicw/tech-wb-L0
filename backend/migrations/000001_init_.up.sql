CREATE TABLE IF NOT EXISTS deliveries (
delivery_id     BIGSERIAL PRIMARY KEY,
name            TEXT,
phone           TEXT,
zip             TEXT,
city            TEXT,
address         TEXT,
region          TEXT,
email           TEXT
);

CREATE TABLE IF NOT EXISTS payments (
payment_id          BIGSERIAL PRIMARY KEY,
transaction         TEXT,
request_id          TEXT,
currency            TEXT,
provider            TEXT,
amount              BIGINT,
payment_dt          BIGINT,
bank                TEXT,
delivery_cost       INT,
goods_total         INT,
custom_fee          INT
);

CREATE TABLE IF NOT EXISTS items (
item_id             BIGSERIAL PRIMARY KEY,
chrt_id             BIGINT,
track_number        TEXT,
price               BIGINT,
rid                 TEXT,
name                TEXT,
sale                INT,
size                TEXT,
total_price         BIGINT,
nm_id               BIGINT,
brand               TEXT,
status              INT
);

CREATE TABLE IF NOT EXISTS orders (
order_uid           TEXT PRIMARY KEY,
track_number        TEXT UNIQUE,
entry               TEXT,
delivery_id         BIGINT NOT NULL REFERENCES deliveries(delivery_id) ON DELETE CASCADE,
payment_id          BIGINT NOT NULL REFERENCES payments(payment_id) ON DELETE CASCADE,
locale              TEXT NOT NULL,
internal_signature  TEXT,
customer_id         TEXT,
delivery_service    TEXT,
shardkey            INT,
sm_id               BIGINT,
date_created        TIMESTAMPTZ,
oof_shard           INT
);

CREATE TABLE IF NOT EXISTS order_items (
order_uid   TEXT NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
item_id     BIGINT NOT NULL REFERENCES items(item_id) ON DELETE CASCADE,
PRIMARY KEY (order_uid, item_id)
);
