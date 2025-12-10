CREATE TABLE products (
  id bigserial PRIMARY KEY,
  sku text UNIQUE NOT NULL,
  name text NOT NULL,
  description text,
  price numeric(12,2) NOT NULL,
  stock integer NOT NULL DEFAULT 0,
  created_at timestamptz DEFAULT now()
);

CREATE TABLE orders (
  id bigserial PRIMARY KEY,
  product_id bigint REFERENCES products(id),
  quantity integer NOT NULL,
  total numeric(12,2) NOT NULL,
  customer_id text,
  created_at timestamptz DEFAULT now()
);

CREATE INDEX idx_orders_customer ON orders(customer_id);
