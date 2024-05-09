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
		// assert.Equal(t, configPath, config.DBFolder)

		assert.NoError(t, err)
		assert.Equal(t, "./nix/nixos-options.nix", config.Sources.NixOS)
		assert.Equal(t, "./nix/hm-options.nix", config.Sources.HomeManager)
		assert.Equal(t, "./nix/darwin-options.nix", config.Sources.Darwin)
	})

	t.Run("Should load config from environment values", func(t *testing.T) {
		os.Setenv("OPTINIX_SOURCES_NIXOS_PATH", "./other/file.nix")
		os.Setenv("OPTINIX_SOURCES_HOME_MANAGER_PATH", "./another/file.nix")
		os.Setenv("OPTINIX_SOURCES_DARWIN_PATH", "./another/darwin.nix")
		os.Setenv("OPTINIX_DB_FOLDER", "/home/test")
		config, err := config.LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, "./other/file.nix", config.Sources.NixOS)
		assert.Equal(t, "./another/file.nix", config.Sources.HomeManager)
		assert.Equal(t, "./another/darwin.nix", config.Sources.Darwin)
		assert.Equal(t, "/home/test", config.DBFolder)
	})
}
