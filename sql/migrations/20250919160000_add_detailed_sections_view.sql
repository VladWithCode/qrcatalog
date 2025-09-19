-- +goose Up
-- +goose StatementBegin

-- View for detailed sections with all related data
CREATE VIEW detailed_sections AS
SELECT
    s.id,
    s.name,
    s.title,
    s.image,
    s.bg_image,
    s.created_at,
    s.updated_at,
    s.search_vector,
    -- Aggregate paragraphs as JSON array, ordered by order field
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                    'id', sp.id,
                    'section_id', sp.section_id,
                    'order', sp.order_idx,
                    'content', sp.content,
                    'created_at', sp.created_at,
                    'updated_at', sp.updated_at
                ) ORDER BY sp.order_idx ASC
            )
            FROM section_paragraphs sp
            WHERE sp.section_id = s.id
        ),
        '[]'::json
    ) as paragraphs,
    -- Aggregate services with their items as JSON array
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                    'id', ss.id,
                    'section_id', ss.section_id,
                    'title', ss.title,
                    'price', ss.price,
                    'description', ss.description,
                    'created_at', ss.created_at,
                    'updated_at', ss.updated_at,
                    'items', COALESCE((
                        SELECT json_agg(
                            json_build_object(
                                'id', ssi.id,
                                'service_id', ssi.service_id,
                                'order', ssi.order_idx,
                                'price', ssi.price,
                                'content', ssi.content,
                                'content_as_list', ssi.content_as_list,
                                'created_at', ssi.created_at,
                                'updated_at', ssi.updated_at
                            ) ORDER BY ssi.order_idx ASC
                        )
                        FROM section_service_items ssi
                        WHERE ssi.service_id = ss.id
                    ), '[]'::json)
                ) ORDER BY ss.created_at ASC
            )
            FROM section_service ss
            WHERE ss.section_id = s.id
        ),
        '[]'::json
    ) as services
FROM sections s
ORDER BY s.created_at DESC;

-- Create indexes to improve view performance
CREATE INDEX idx_section_paragraphs_section_id ON section_paragraphs(section_id);
CREATE INDEX idx_section_service_section_id ON section_service(section_id);
CREATE INDEX idx_section_service_items_service_id ON section_service_items(service_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_section_service_items_service_id;
DROP INDEX IF EXISTS idx_section_service_section_id;
DROP INDEX IF EXISTS idx_section_paragraphs_section_id;

DROP VIEW IF EXISTS detailed_sections;

-- +goose StatementEnd
