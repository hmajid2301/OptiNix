-- name: FindOptions :many
SELECT
    options.id,
    options.name,
    options.description,
    options.type,
    options.default_value,
    options.example,
    GROUP_CONCAT(sources.url) AS source_list
FROM
    options
LEFT JOIN
    source_options ON options.id = source_options.option_id
LEFT JOIN
    sources ON source_options.source_id = sources.id
WHERE
	name LIKE ?
GROUP BY
    options.id
LIMIT
	10;

-- name: AddOption :one
INSERT INTO options (name, description, type, default_value, example) VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: AddSource :one
INSERT INTO sources (url) VALUES (?) RETURNING *;

-- name: AddSourceOption :one
INSERT INTO source_options (source_id, option_id) VALUES (?, ?) RETURNING *;

-- name: GetLastUpdated :one
SELECT
    options.updated_at
FROM
    options
ORDER BY 
    options.updated_at DESC
LIMIT
	1;
