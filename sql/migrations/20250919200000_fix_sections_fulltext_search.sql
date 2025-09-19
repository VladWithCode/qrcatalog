-- +goose Up
-- +goose StatementBegin

-- Drop the existing broken function and trigger
DROP TRIGGER IF EXISTS update_sections_search_vector ON sections;
DROP FUNCTION IF EXISTS update_section_search_vector();

-- Create corrected function to update section search vector
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

-- Recreate trigger to automatically update search vector
CREATE TRIGGER update_sections_search_vector
    BEFORE INSERT OR UPDATE ON sections
    FOR EACH ROW EXECUTE FUNCTION update_section_search_vector();

-- Update all existing records to populate search vectors correctly
UPDATE sections SET updated_at = CURRENT_TIMESTAMP;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop trigger and function
DROP TRIGGER IF EXISTS update_sections_search_vector ON sections;
DROP FUNCTION IF EXISTS update_section_search_vector();

-- +goose StatementEnd