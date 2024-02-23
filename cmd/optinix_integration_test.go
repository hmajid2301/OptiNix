package cmd_test

import (
	"testing"
)

func TestRootCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// root := &cobra.Command{RunE: cmd.FindOptions}
	//
	// out, err := execute(t, root, "")
	// assert.NoError(t, err)
	// assert.Contains(t, out, "")
}

// func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
// 	t.Helper()
//
// 	buf := new(bytes.Buffer)
// 	c.SetOut(buf)
// 	c.SetErr(buf)
// 	c.SetArgs(args)
//
// 	err := c.Execute()
// 	return buf.String(), err
// }
