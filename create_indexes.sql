CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS idx_taxes_business_id ON product_taxes(business_id);
CREATE INDEX IF NOT EXISTS idx_categories_business_id ON product_categories(business_id);
CREATE INDEX IF NOT EXISTS idx_discounts_business_id ON product_discounts(business_id);
CREATE INDEX IF NOT EXISTS idx_units_business_id ON product_units(business_id);
CREATE INDEX IF NOT EXISTS idx_brands_business_id ON brands(business_id);
CREATE INDEX IF NOT EXISTS idx_bundles_business_id ON product_bundles(business_id);
CREATE INDEX IF NOT EXISTS idx_products_business_id ON products(business_id);

CREATE INDEX IF NOT EXISTS idx_taxes_biz_name ON product_taxes(business_id, name);
CREATE INDEX IF NOT EXISTS idx_categories_biz_name ON product_categories(business_id, name);
CREATE INDEX IF NOT EXISTS idx_discounts_biz_name ON product_discounts(business_id, name);
CREATE INDEX IF NOT EXISTS idx_units_biz_name ON product_units(business_id, name);
CREATE INDEX IF NOT EXISTS idx_brands_biz_name ON brands(business_id, name);
CREATE INDEX IF NOT EXISTS idx_bundles_biz_name ON product_bundles(business_id, name);
CREATE INDEX IF NOT EXISTS idx_products_biz_name ON products(business_id, name);

CREATE INDEX IF NOT EXISTS idx_taxes_name_trgm ON product_taxes USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_categories_name_trgm ON product_categories USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_discounts_name_trgm ON product_discounts USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_units_name_trgm ON product_units USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_brands_name_trgm ON brands USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_bundles_name_trgm ON product_bundles USING gin (name gin_trgm_ops);
