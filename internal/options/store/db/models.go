// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package sqlc

import (
	"database/sql"
)

type Option struct {
	ID           int64
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	OptionName   string
	Description  string
	OptionType   string
	OptionFrom   string
	DefaultValue string
	Example      string
}

type OptionsFt struct {
	OptionID    string
	OptionName  string
	Description string
}

type Source struct {
	ID        int64
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	Url       string
}

type SourceOption struct {
	SourceID  int64
	OptionID  int64
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}
