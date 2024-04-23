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

		assert.NoError(t, err)
		assert.Equal(t, "https://nixos.org/manual/nixos/unstable/options", config.Sources.NixOSURL)
		assert.Equal(t, "https://nix-community.github.io/home-manager/options.xhtml", config.Sources.HomeManagerURL)
		assert.Equal(t, 3, config.Retries)
		assert.Equal(t, 30, config.Timeout)
	})

	t.Run("Should load config from environment values", func(t *testing.T) {
		os.Setenv("NIXOS_URL", "http://docker:8080/manual/nixos/unstable/options")
		os.Setenv("HOME_MANAGER_URL", "http://docker:8080/home-manager/options.xhtml")
		config, err := config.LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, "http://docker:8080/manual/nixos/unstable/options", config.Sources.NixOSURL)
		assert.Equal(t, "http://docker:8080/home-manager/options.xhtml", config.Sources.HomeManagerURL)
		assert.Equal(t, 3, config.Retries)
		assert.Equal(t, 30, config.Timeout)
	})
}
