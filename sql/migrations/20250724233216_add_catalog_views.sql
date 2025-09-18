-- +goose Up
-- +goose StatementBegin

-- View for catalog categories with product counts
CREATE VIEW catalog_categories AS
SELECT 
    c.id,
    c.name,
    COUNT(p.id) as product_count
FROM public.categories c
LEFT JOIN public.products p ON c.id = p.category_id
GROUP BY c.id, c.name
ORDER BY c.name;

-- View for catalog products with all related data
CREATE VIEW catalog_products AS
SELECT 
    p.id,
    p.name,
    p.description,
    p.long_description,
    p.slug,
    p.category_id,
    c.name as category_name,
    COALESCE(main_img.filename, '') as image_url,
    p.price,
    p.unit,
    p.available,
    p.quantity,
    p.search_vector,
    -- Aggregate gallery images as JSON array
    COALESCE(
        (
            SELECT json_agg(i.filename ORDER BY i.filename)
            FROM public.images_products ip
            JOIN public.images i ON ip.image_id = i.id
            WHERE ip.product_id = p.id
        ),
        '[]'::json
    ) as images
FROM public.products p
LEFT JOIN public.categories c ON p.category_id = c.id
LEFT JOIN public.images main_img ON p.main_img_id = main_img.id
ORDER BY p.name;

CREATE INDEX idx_images_products_product_id ON public.images_products(product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_products_category;
DROP INDEX IF EXISTS idx_images_products_product_id;

DROP VIEW IF EXISTS catalog_products;
DROP VIEW IF EXISTS catalog_categories;
-- +goose StatementEnd
