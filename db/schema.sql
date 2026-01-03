CREATE TABLE IF NOT EXISTS options (
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    option_name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    option_type TEXT NOT NULL,
    option_from TEXT NOT NULL,
    default_value TEXT NOT NULL,
    example TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sources (
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    url TEXT NOT NULL
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

CREATE VIRTUAL TABLE IF NOT EXISTS options_fts USING fts5 (
    option_id, option_name, description
);

CREATE TRIGGER IF NOT EXISTS insert_options_fts
AFTER INSERT ON options
BEGIN
    INSERT INTO options_fts (option_id, option_name, description)
    VALUES (new.id, new.option_name, new.description);
END;

CREATE TRIGGER IF NOT EXISTS update_options_fts
AFTER UPDATE ON options
FOR EACH ROW
BEGIN
    UPDATE options_fts SET
        option_name = new.option_name, description = new.description
    WHERE option_id = old.id;
END;

CREATE TRIGGER IF NOT EXISTS delete_options_fts
AFTER DELETE ON options
FOR EACH ROW
BEGIN
    DELETE FROM options_fts
    WHERE option_id = old.id;
END;

-- CREATE INDEX IF NOT EXISTS options_name_idx ON options (option_name);
CREATE UNIQUE INDEX IF NOT EXISTS sources_url_idx ON sources (url);
