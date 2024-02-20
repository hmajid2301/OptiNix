package cmd_test

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"gitlab.com/majiy00/go/clis/optinix/cmd"
)

func TestRootCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	root := &cobra.Command{RunE: cmd.FindOptions}

	out, err := execute(t, root, "alacritty")
	assert.NoError(t, err)
	assert.Contains(t, out, "Whether to enable Alacritty.")
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
