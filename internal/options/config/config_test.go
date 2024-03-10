package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options/config"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Should load config with default values", func(t *testing.T) {
		config, err := config.LoadConfig()

		assert.Nil(t, err)
		assert.Equal(t, "https://nixos.org/manual/nixos/unstable/options", config.Sources.NixOSURL)
		assert.Equal(t, "https://nix-community.github.io/home-manager/options.xhtml", config.Sources.HomeManagerURL)
		assert.Equal(t, 3, config.Retries)
		assert.Equal(t, 30, config.Timeout)
	})
}
