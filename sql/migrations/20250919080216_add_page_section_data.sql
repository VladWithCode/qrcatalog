-- +goose Up
-- +goose StatementBegin
CREATE TABLE sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) UNIQUE NOT NULL,
    title TEXT,
    image VARCHAR(256),
    bg_image VARCHAR(256),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE section_paragraphs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id UUID NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    order_idx INT NOT NULL,
    content TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE section_service (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id UUID NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    -- Price is multiplied by 100 (e.g. 1000 = 10$)
    price INT,
    -- description may be empty and a list of items will be shown instead
    -- (e.g. a list of included services/products)
    -- If description is not empty, it will be shown as a paragraph, and the list of items will be shown below
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE section_service_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES section_service(id) ON DELETE CASCADE,
    order_idx INT NOT NULL DEFAULT 0,
    price INT NOT NULL,
    content TEXT NOT NULL,
    content_as_list BOOLEAN DEFAULT FALSE NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create triggers to automatically update updated_at column
CREATE TRIGGER update_sections_updated_at
    BEFORE UPDATE ON sections
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_section_paragraphs_updated_at
    BEFORE UPDATE ON section_paragraphs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_section_service_updated_at
    BEFORE UPDATE ON section_service
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_section_service_items_updated_at
    BEFORE UPDATE ON section_service_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop triggers
DROP TRIGGER IF EXISTS update_section_service_items_updated_at ON section_service_items;
DROP TRIGGER IF EXISTS update_section_service_updated_at ON section_service;
DROP TRIGGER IF EXISTS update_section_paragraphs_updated_at ON section_paragraphs;
DROP TRIGGER IF EXISTS update_sections_updated_at ON sections;

-- Drop tables
DROP TABLE section_service_items;
DROP TABLE section_service;
DROP TABLE section_paragraphs;
DROP TABLE sections;
-- +goose StatementEnd
