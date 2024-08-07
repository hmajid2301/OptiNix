-- name: FindOptions :many
SELECT
    o.id,
    o.option_name,
    o.description,
    o.option_type,
    o.default_value,
    o.example,
    o.option_from,
    IFNULL(GROUP_CONCAT(s.url), '') AS source_list
FROM
    options o
LEFT JOIN
    source_options so ON o.id = so.option_id
LEFT JOIN
    sources s ON so.source_id = s.id
WHERE
    o.id IN (
        SELECT option_id FROM options_fts WHERE options_fts.option_name MATCH ?
    )
GROUP BY
    o.id
LIMIT
    ?;

-- name: GetAllOptions :many
SELECT
    o.id,
    o.option_name,
    o.description,
    o.option_type,
    o.default_value,
    o.example,
    o.option_from,
    IFNULL(GROUP_CONCAT(s.url), '') AS source_list
FROM
    options o
LEFT JOIN
    source_options so ON o.id = so.option_id
LEFT JOIN
    sources s ON so.source_id = s.id
GROUP BY
    o.id;

-- name: AddOption :one
INSERT OR REPLACE INTO
   options (option_name, description, option_type, option_from, default_value, example)
VALUES (?, ?, ?, ?, ?, ?) RETURNING *;

-- name: AddSource :one
INSERT INTO sources (url) VALUES (?) ON CONFLICT(url) DO UPDATE SET url = excluded.url RETURNING *;

-- name: AddSourceOption :one
INSERT INTO source_options (source_id, option_id) VALUES (?, ?) RETURNING *;

-- name: GetOptionCount :one
SELECT COUNT(*) FROM options;
