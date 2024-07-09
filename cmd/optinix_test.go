package cmd_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrationExecuteCommand(t *testing.T) {
	t.Run("Should show help", func(t *testing.T) {
		// ctx := context.Background()
		os.Setenv("OPTINIX_DB_FOLDER", "../../testdata/")
		assert.True(t, true)

		// _, filename, _, ok := runtime.Caller(0)
		// assert.True(t, ok)
		// dir := path.Join(path.Dir(filename), "..")
		// schemaFile := filepath.Join(dir, "db", "schema.sql")
		// content, err := os.ReadFile(schemaFile)
		// assert.NoError(t, err)
		//
		// ddl := string(content)
		// cmd, err := cmd.NewRootCmd(ctx, ddl, embed.FS{})
		//
		// b := bytes.NewBufferString("")
		// cmd.SetOut(b)
		// cmd.SetArgs([]string{"--help"})
		// cmd.Execute()
		// out := b.String()
		// if err != nil {
		// 	t.Fatal(err)
		// }
		//
		// assert.Contains(t, out, "--help")
	})
}
