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

	"gitlab.com/hmajid2301/optinix/internal/options"
)

func TestIntegrationCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	os.Setenv("OPTINIX_SOURCES_NIXOS_PATH", "./nix/nixos-options.nix")
	os.Setenv("OPTINIX_SOURCES_HOME_MANAGER_PATH", "./nix/hm-options.nix")
	os.Setenv("OPTINIX_SOURCES_DARWIN_PATH", "./nix/darwin-options.nix")
	os.Setenv("OPTINIX_DB_FOLDER", "../testdata")

	db, err := options.GetDB("../testdata")
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

	// TODO: fix this test to use bubbletea
	root := &cobra.Command{RunE: func(_ *cobra.Command, _ []string) error {
		// return cmd.FindOptions(ctx, db, args[0])
		return nil
	},
	}
	// out, err := execute(t, root, "appstream")
	_, err = execute(t, root, "appstream")
	assert.NoError(t, err)
	// assert.Contains(t, out, "appstream")

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
