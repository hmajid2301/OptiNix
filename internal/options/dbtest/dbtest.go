package dbtest

import (
	"context"
	"database/sql"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	// used to connect to sqlite
	_ "modernc.org/sqlite"
)

func CreateDB(ctx context.Context, t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	assert.NoError(t, err)
	assert.NoError(t, err)
	_, filename, _, ok := runtime.Caller(0)
	assert.True(t, ok)
	dir := path.Join(path.Dir(filename), "..")
	schemaFile := filepath.Join(dir, "../", "../", "db", "schema.sql")
	content, err := os.ReadFile(schemaFile)
	assert.NoError(t, err)
	ddl := string(content)
	_, err = db.ExecContext(ctx, ddl)
	assert.NoError(t, err)
	return db
}
