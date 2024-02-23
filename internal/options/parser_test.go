//go:build unit
// +build unit

package options_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"gitlab.com/hmajid2301/optinix/internal/options"
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
				Sources: []string{
					"https://github.com/NixOS/nixpkgs/blob/release-23.11/nixos/modules/config/appstream.nix",
				},
			},
			{
				Name: "boot.enableContainers",
				Description: "Whether to enable support for NixOS containers. Defaults to true " +
					"(at no cost if containers are not actually used).",
				Type:    "boolean",
				Default: "true",
				Sources: []string{
					"https://github.com/NixOS/nixpkgs/blob/release-23.11/nixos/modules/virtualisation/nixos-containers.nix",
				},
			},
		}

		assert.Equal(t, expectedOpts, opts)
	})

	t.Run("Should parse empty NixOS options", func(t *testing.T) {
		data, err := os.ReadFile("../../testdata/nixos_options_empty.html")
		assert.NoError(t, err)
		reader := bytes.NewReader(data)

		node, err := html.Parse(reader)
		assert.NoError(t, err)
		opts := options.Parse(node)

		expectedOpts := []options.Option{}
		assert.Equal(t, expectedOpts, opts)
	})

	t.Run("Should parse Home Manager options", func(t *testing.T) {
		data, err := os.ReadFile("../../testdata/home_manager_options.html")
		assert.NoError(t, err)
		reader := bytes.NewReader(data)

		node, err := html.Parse(reader)
		assert.NoError(t, err)
		opts := options.Parse(node)

		expectedOpts := []options.Option{
			{
				Name:        "accounts.calendar.accounts",
				Description: "List of calendars.",
				Type:        "attribute set of (submodule)",
				Default:     "{ }",
				Sources: []string{
					"https://github.com/nix-community/home-manager/blob/master/modules/programs/qcal.nix",
					"https://github.com/nix-community/home-manager/blob/master/modules/accounts/calendar.nix",
				},
			},
			{
				Name:        "accounts.calendar.accounts.<name>.khal.enable",
				Description: "Whether to enable khal access.",
				Type:        "boolean",
				Default:     "false",
				Example:     "true",
				Sources: []string{
					"https://github.com/nix-community/home-manager/blob/master/modules/accounts/calendar.nix",
				},
			},
		}

		assert.Equal(t, expectedOpts, opts)
	})
}
