package cmd_test

import (
	"bytes"
	"context"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/cmd"
	"gitlab.com/hmajid2301/optinix/internal/options/optionstest"
)

func TestIntegrationCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	os.Setenv("OPTINIX_SOURCES_NIXOS_URL", optionstest.GetHost("/manual/nixos/unstable/options"))
	os.Setenv("OPTINIX_SOURCES_HOME_MANAGER_URL", optionstest.GetHost("/home-manager/options.xhtml"))
	os.Setenv("OPTINIX_DB_FOLDER", "../testdata")

	db, err := cmd.GetDB("../testdata")
	assert.NoError(t, err)

	_, filename, _, ok := runtime.Caller(0)
	assert.True(t, ok)
	dir := path.Join(path.Dir(filename), "..")
	schemaFile := filepath.Join(dir, "db", "schema.sql")
	content, err := os.ReadFile(schemaFile)
	assert.NoError(t, err)

	ctx := context.TODO()
	ddl := string(content)
	_, err = db.ExecContext(ctx, ddl)
	assert.NoError(t, err)

	root := &cobra.Command{RunE: cmd.FindOptions}
	out, err := execute(t, root, "appstream")
	assert.NoError(t, err)
	assert.Contains(t, out, "appstream")

	t.Cleanup(func() {
		os.Remove("../testdata/optinix.db")
		os.Remove("../testdata/optinix.db-shm")
		os.Remove("../testdata/optinix.db-wal")
	})
}

func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()
	return buf.String(), err
}
