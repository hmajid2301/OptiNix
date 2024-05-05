package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options/config"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Should load config with default values", func(t *testing.T) {
		config, err := config.LoadConfig()
		// TODO:: how to test db folder
		// state := os.Getenv("XDG_DATA_HOME")
		// configPath := filepath.Join(state, "optinix")

		assert.NoError(t, err)
		assert.Equal(t, "https://nixos.org/manual/nixos/unstable/options", config.Sources.NixOSURL)
		assert.Equal(t, "https://nix-community.github.io/home-manager/options.xhtml", config.Sources.HomeManagerURL)
		assert.Equal(t, 3, config.Retries)
		assert.Equal(t, 30, config.Timeout)
		// assert.Equal(t, configPath, config.DBFolder)
	})

	t.Run("Should load config from environment values", func(t *testing.T) {
		os.Setenv("OPTINIX_SOURCES_NIXOS_URL", "http://docker:8080/manual/nixos/unstable/options")
		os.Setenv("OPTINIX_SOURCES_HOME_MANAGER_URL", "http://docker:8080/home-manager/options.xhtml")
		os.Setenv("OPTINIX_DB_FOLDER", "/home/test")
		config, err := config.LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, "http://docker:8080/manual/nixos/unstable/options", config.Sources.NixOSURL)
		assert.Equal(t, "http://docker:8080/home-manager/options.xhtml", config.Sources.HomeManagerURL)
		assert.Equal(t, 3, config.Retries)
		assert.Equal(t, 30, config.Timeout)
		assert.Equal(t, "/home/test", config.DBFolder)
	})
}
