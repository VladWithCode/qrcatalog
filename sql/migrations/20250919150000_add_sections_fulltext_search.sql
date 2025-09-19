-- +goose Up
-- +goose StatementBegin

-- Add tsvector column for full-text search on sections
ALTER TABLE sections ADD COLUMN search_vector tsvector;

-- Create function to update section search vector
CREATE OR REPLACE FUNCTION update_section_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    -- Combine section data with weighted importance for search
    NEW.search_vector :=
        -- Section name and title (highest weight - most important)
        setweight(to_tsvector('spanish', COALESCE(NEW.name, '')), 'A') ||
        setweight(to_tsvector('spanish', COALESCE(NEW.title, '')), 'A') ||

        -- Paragraph content (medium-high weight)
        setweight(to_tsvector('spanish', COALESCE(
            (SELECT string_agg(sp.content, ' ')
             FROM section_paragraphs sp
             WHERE sp.section_id = NEW.id), ''
        )), 'B') ||

        -- Service titles and descriptions (medium weight)
        setweight(to_tsvector('spanish', COALESCE(
            (SELECT string_agg(
                CONCAT(ss.title, ' ', COALESCE(ss.description, '')), ' '
             ) FROM section_service ss
             WHERE ss.section_id = NEW.id), ''
        )), 'B') ||

        -- Service item content (lower weight - least specific)
        setweight(to_tsvector('spanish', COALESCE(
            (SELECT string_agg(ssi.content, ' ')
             FROM section_service ss
             JOIN section_service_items ssi ON ss.id = ssi.service_id
             WHERE ss.section_id = NEW.id), ''
        )), 'C');

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update search vector
CREATE TRIGGER update_sections_search_vector
    BEFORE INSERT OR UPDATE ON sections
    FOR EACH ROW EXECUTE FUNCTION update_section_search_vector();

-- Update existing records to populate search vectors
-- Force trigger execution for all existing sections
UPDATE sections SET updated_at = updated_at WHERE id IN (SELECT id FROM sections);

-- Create GIN index for fast full-text search
CREATE INDEX idx_sections_search_vector ON sections USING gin(search_vector);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop index
DROP INDEX IF EXISTS idx_sections_search_vector;

-- Drop trigger
DROP TRIGGER IF EXISTS update_sections_search_vector ON sections;

-- Drop function
DROP FUNCTION IF EXISTS update_section_search_vector();

-- Drop column
ALTER TABLE sections DROP COLUMN IF EXISTS search_vector;

-- +goose StatementEnd