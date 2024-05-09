package options_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options"
)

func TestParse(t *testing.T) {
	t.Run("Should successfully parse options from HM JSON options file", func(t *testing.T) {
		content, err := os.ReadFile("../../testdata/hm-options.json")
		assert.NoError(t, err)

		opts, err := options.ParseOptions(content)
		assert.Len(t, opts, 97)
		assert.NoError(t, err)
		option := options.Option{
			Name:        "accounts.calendar.accounts.<name>.vdirsyncer.enable",
			Description: "Whether to enable synchronization using vdirsyncer.",
			Type:        "boolean",
			Default:     "false",
			Example:     "true",
			Sources: []string{
				"https://github.com/nix-community/home-manager/blob/master/modules/accounts/calendar.nix",
			},
		}

		assert.Contains(t, opts, option)
	})

	t.Run("Should successfully parse options from Darwin JSON options file", func(t *testing.T) {
		content, err := os.ReadFile("../../testdata/darwin-options.json")
		assert.NoError(t, err)

		opts, err := options.ParseOptions(content)
		assert.NoError(t, err)
		assert.Len(t, opts, 72)
		option := options.Option{
			Name:        "homebrew.brews.*.conflicts_with",
			Description: "List of formulae that should be unlinked and their services stopped (if they are\ninstalled).\n",
			Type:        "null or (list of string)",
			Default:     "null",
			Example:     "",
			Sources: []string{
				"https://github.com/LnL7/nix-darwin/blob/master/modules/homebrew.nix",
			},
		}

		assert.Contains(t, opts, option)
	})

	t.Run("Should successfully parse options from NixOS JSON options file", func(t *testing.T) {
		content, err := os.ReadFile("../../testdata/nixos-options.json")
		assert.NoError(t, err)

		opts, err := options.ParseOptions(content)
		assert.NoError(t, err)
		assert.Len(t, opts, 97)
		option := options.Option{
			Name:        "accounts.calendar.accounts.<name>.vdirsyncer.enable",
			Description: "Whether to enable synchronization using vdirsyncer.",
			Type:        "boolean",
			Default:     "false",
			Example:     "true",
			Sources: []string{
				"https://github.com/nix-community/home-manager/blob/master/modules/accounts/calendar.nix",
			},
		}

		assert.Contains(t, opts, option)
	})

	t.Run("Should successfully parse options empty options file", func(t *testing.T) {
		content, err := os.ReadFile("../../testdata/empty-options.json")
		assert.NoError(t, err)

		opts, err := options.ParseOptions(content)
		assert.NoError(t, err)
		assert.Len(t, opts, 0)
	})

	t.Run("Should fail to parse random JSON file", func(t *testing.T) {
		content, err := os.ReadFile("../../testdata/wrong-file.json")
		assert.NoError(t, err)

		_, err = options.ParseOptions(content)
		assert.ErrorContains(t, err, "cannot unmarshal array into Go value of type map[string]options.OptionFile")
	})

	t.Run("Should fail parse invalid json options file", func(t *testing.T) {
		content, err := os.ReadFile("../../testdata/invalid-json-options.json")
		assert.NoError(t, err)

		_, err = options.ParseOptions(content)
		assert.ErrorContains(t, err, "unexpected end of JSON input")
	})
}
