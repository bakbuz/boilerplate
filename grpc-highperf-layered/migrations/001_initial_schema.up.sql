CREATE TABLE IF NOT EXISTS catalog.brands
(
    id SERIAL NOT NULL,
    name character varying(100) NOT NULL,
    slug character varying(100),
    logo character varying(100),
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS catalog.products
(
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    brand_id integer NOT NULL,
    name character varying(100) NOT NULL,
    sku character varying(100) ,
    summary character varying(100) ,
    storyline character varying(1000) ,
    stock_quantity integer NOT NULL DEFAULT 0,
    price money,
    deleted boolean NOT NULL DEFAULT false,
	created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by UUID,
    updated_at TIMESTAMPTZ,
	deleted_by UUID,
    deleted_at TIMESTAMPTZ
);

-- Optimize for Supabase
ALTER TABLE catalog.brands ENABLE ROW LEVEL SECURITY;
ALTER TABLE catalog.products ENABLE ROW LEVEL SECURITY;