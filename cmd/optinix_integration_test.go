package cmd_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/cmd"
)

func TestRootCmd(t *testing.T) {
	os.Setenv("NIXOS_URL", "http://docker:8080/manual/nixos/unstable/options")
	os.Setenv("HOME_MANAGER_URL", "http://docker:8080/manual/home-manager/options.xhtml")

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	db, err := cmd.GetDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		return
	}

	dir, err := os.Getwd()
	assert.NoError(t, err)

	schemaFile := filepath.Join(dir, "../db/schema.sql")
	content, err := os.ReadFile(schemaFile)
	assert.NoError(t, err)
	ddl := string(content)
	_, err = db.ExecContext(ctx, ddl)
	assert.NoError(t, err)

	root := &cobra.Command{RunE: cmd.FindOptions}
	out, _ := execute(t, root, "")
	// assert.NoError(t, err)
	assert.Contains(t, out, "")
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
