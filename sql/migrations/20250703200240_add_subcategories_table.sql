-- +goose Up
-- +goose StatementBegin
CREATE TABLE subcategories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    description VARCHAR(360) NOT NULL,
    long_description VARCHAR(512),
    display_img UUID REFERENCES images(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE subcategories;
-- +goose StatementEnd
