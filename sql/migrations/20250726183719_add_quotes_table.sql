-- +goose Up
-- +goose StatementBegin
CREATE TABLE quotes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_name VARCHAR(256) NOT NULL,
    customer_phone VARCHAR(256) NOT NULL,
    customer_email VARCHAR(256) NOT NULL,
    cart_id UUID REFERENCES carts(id) ON DELETE SET NULL,
    request_type VARCHAR(64) NOT NULL, -- contacto, cotización
    status VARCHAR(64) NOT NULL DEFAULT 'pendiente', -- pendiente, procesada, en_progreso, cancelada
    comments TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER quotes_set_updated_at BEFORE UPDATE ON quotes
    FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER quotes_set_updated_at ON quotes;

DROP TABLE quotes;
-- +goose StatementEnd
