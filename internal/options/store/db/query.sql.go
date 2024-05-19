// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"
)

const addOption = `-- name: AddOption :one
INSERT INTO options (option_name, description, option_type, default_value, example) VALUES (?, ?, ?, ?, ?) RETURNING id, created_at, updated_at, option_name, description, option_type, default_value, example
`

type AddOptionParams struct {
	OptionName   string
	Description  string
	OptionType   string
	DefaultValue string
	Example      string
}

func (q *Queries) AddOption(ctx context.Context, arg AddOptionParams) (Option, error) {
	row := q.db.QueryRowContext(ctx, addOption,
		arg.OptionName,
		arg.Description,
		arg.OptionType,
		arg.DefaultValue,
		arg.Example,
	)
	var i Option
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OptionName,
		&i.Description,
		&i.OptionType,
		&i.DefaultValue,
		&i.Example,
	)
	return i, err
}

const addSource = `-- name: AddSource :one
INSERT INTO sources (url) VALUES (?) ON CONFLICT(url) DO UPDATE SET url = excluded.url RETURNING id, created_at, updated_at, url
`

func (q *Queries) AddSource(ctx context.Context, url string) (Source, error) {
	row := q.db.QueryRowContext(ctx, addSource, url)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Url,
	)
	return i, err
}

const addSourceOption = `-- name: AddSourceOption :one
INSERT INTO source_options (source_id, option_id) VALUES (?, ?) RETURNING source_id, option_id, created_at, updated_at
`

type AddSourceOptionParams struct {
	SourceID int64
	OptionID int64
}

func (q *Queries) AddSourceOption(ctx context.Context, arg AddSourceOptionParams) (SourceOption, error) {
	row := q.db.QueryRowContext(ctx, addSourceOption, arg.SourceID, arg.OptionID)
	var i SourceOption
	err := row.Scan(
		&i.SourceID,
		&i.OptionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findOptions = `-- name: FindOptions :many
SELECT
    o.id,
    o.option_name,
    o.description,
    o.option_type,
    o.default_value,
    o.example,
    GROUP_CONCAT(s.url) AS source_list
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
    ?
`

type FindOptionsParams struct {
	OptionName string
	Limit      int64
}

type FindOptionsRow struct {
	ID           int64
	OptionName   string
	Description  string
	OptionType   string
	DefaultValue string
	Example      string
	SourceList   string
}

func (q *Queries) FindOptions(ctx context.Context, arg FindOptionsParams) ([]FindOptionsRow, error) {
	rows, err := q.db.QueryContext(ctx, findOptions, arg.OptionName, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindOptionsRow
	for rows.Next() {
		var i FindOptionsRow
		if err := rows.Scan(
			&i.ID,
			&i.OptionName,
			&i.Description,
			&i.OptionType,
			&i.DefaultValue,
			&i.Example,
			&i.SourceList,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastUpdated = `-- name: GetLastUpdated :one
SELECT
    options.updated_at
FROM
    options
ORDER BY
    options.updated_at DESC
LIMIT
	1
`

func (q *Queries) GetLastUpdated(ctx context.Context) (sql.NullTime, error) {
	row := q.db.QueryRowContext(ctx, getLastUpdated)
	var updated_at sql.NullTime
	err := row.Scan(&updated_at)
	return updated_at, err
}
