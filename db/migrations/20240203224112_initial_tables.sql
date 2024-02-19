-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sources (
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    url TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS options (
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    type TEXT NOT NULL,
    default_value TEXT NOT NULL,
    example TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS source_options (
    source_id INTEGER NOT NULL,
    option_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (source_id, option_id),
    FOREIGN KEY (source_id) REFERENCES sources (id),
    FOREIGN KEY (option_id) REFERENCES options (id)
);

CREATE INDEX IF NOT EXISTS options_name ON options (name);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX options_name;

DROP TABLE source_options;

DROP TABLE sources;
DROP TABLE options;

-- +goose StatementEnd
