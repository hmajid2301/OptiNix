package options_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/hmajid2301/opnix/internal/options"
	"golang.org/x/net/html"
)

func TestParse(t *testing.T) {
	t.Run("Should parse NixOS options", func(t *testing.T) {
		data, err := os.ReadFile("../../testdata/nixos_options.html")
		assert.NoError(t, err)
		reader := bytes.NewReader(data)

		node, err := html.Parse(reader)
		assert.NoError(t, err)
		opts := options.Parse(node)

		expectedOpts := []options.Option{
			{
				Name: "appstream.enable",
				Description: "Whether to install files to support the " +
					"[AppStream metadata specification](https://www.freedesktop.org/software/appstream/docs/index.html).",
				Type:    "boolean",
				Default: "true",
				Source:  "https://github.com/NixOS/nixpkgs/blob/release-23.11/nixos/modules/config/appstream.nix",
			},
			{
				Name: "boot.enableContainers",
				Description: "Whether to enable support for NixOS containers. Defaults to true " +
					"(at no cost if containers are not actually used).",
				Type:    "boolean",
				Default: "true",
				Source:  "https://github.com/NixOS/nixpkgs/blob/release-23.11/nixos/modules/virtualisation/nixos-containers.nix",
			},
		}

		assert.Equal(t, expectedOpts, opts)
	})
}
