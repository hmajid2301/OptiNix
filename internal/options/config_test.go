package options_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Should load config with default values", func(t *testing.T) {
		_, err := options.LoadConfig()
		// TODO:: how to test db folder
		// state := os.Getenv("XDG_DATA_HOME")
		// configPath := filepath.Join(state, "optinix")
		// assert.Equal(t, configPath, config.DBFolder)

		assert.NoError(t, err)
	})

	t.Run("Should load config from environment values", func(t *testing.T) {
		os.Setenv("OPTINIX_DB_FOLDER", "/home/test")
		config, err := options.LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, "/home/test", config.DBFolder)
	})
}
