package tui_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/stretchr/testify/assert"

	"gitlab.com/hmajid2301/optinix/internal/options/outputs/tui"
)

func TestIntegrationTUI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	t.Run("Should successfully show options", func(t *testing.T) {
		getOptionsFunc := func() tea.Msg {
			return tui.DoneMsg{
				List: []list.Item{
					tui.Item{
						OptionName:   "accounts.calendar.accounts.<name>.vdirsyncer.enable",
						Desc:         "Whether to enable synchronization using vdirsyncer.",
						OptionType:   "boolean",
						DefaultValue: "false",
						Example:      "true",
						Sources: []string{
							"https://github.com/nix-community/home-manager/blob/master/modules/accounts/calendar.nix",
						},
						OptionFrom: "Home Manager",
					},
					tui.Item{
						OptionName:   "accounts.calendar.accounts.<name>.vdirsyncer.package",
						Desc:         "vdirsyncer package to use.",
						OptionType:   "package",
						DefaultValue: "pkgs.vdirsyncer",
						Sources: []string{
							"https://github.com/nix-community/home-manager/blob/master/modules/programs/vdirsyncer.nix",
						},
						OptionFrom: "Home Manager",
					},
					tui.Item{
						OptionName:   "accounts.calendar.accounts.<name>.vdirsyncer.statusPath",
						Desc:         "vdirsyncer ...",
						OptionType:   "string",
						DefaultValue: "$XDG_DATA_HOME/vdirsyncer/status",
						Sources: []string{
							"https://github.com/nix-community/home-manager/blob/master/modules/programs/vdirsyncer.nix",
						},
						OptionFrom: "Home Manager",
					},
				},
			}
		}

		model, err := tui.NewTUI(getOptionsFunc)
		assert.NoError(t, err)
		tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(300, 100))

		teatest.WaitFor(
			t, tm.Output(),
			func(bts []byte) bool {
				return bytes.Contains(bts, []byte("accounts.calendar.accounts.<name>.vdirsyncer.enable"))
			},
			teatest.WithCheckInterval(time.Millisecond*100),
			teatest.WithDuration(time.Second*3),
		)

		tm.Send(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune("q"),
		})
	})

	t.Run("Should successfully show detailed option view", func(t *testing.T) {
		getOptionsFunc := func() tea.Msg {
			return tui.DoneMsg{
				List: []list.Item{
					tui.Item{
						OptionName:   "accounts.calendar.accounts.<name>.vdirsyncer.enable",
						Desc:         "Whether to enable synchronization using vdirsyncer.",
						OptionType:   "boolean",
						DefaultValue: "false",
						Example:      "true",
						Sources: []string{
							"https://github.com/nix-community/home-manager/blob/master/modules/accounts/calendar.nix",
						},
						OptionFrom: "Home Manager",
					},
					tui.Item{
						OptionName:   "accounts.calendar.accounts.<name>.vdirsyncer.package",
						Desc:         "vdirsyncer package to use.",
						OptionType:   "package",
						DefaultValue: "pkgs.vdirsyncer",
						Sources: []string{
							"https://github.com/nix-community/home-manager/blob/master/modules/programs/vdirsyncer.nix",
						},
						OptionFrom: "Home Manager",
					},
					tui.Item{
						OptionName:   "accounts.calendar.accounts.<name>.vdirsyncer.statusPath",
						Desc:         "vdirsyncer ...",
						OptionType:   "string",
						DefaultValue: "$XDG_DATA_HOME/vdirsyncer/status",
						Sources: []string{
							"https://github.com/nix-community/home-manager/blob/master/modules/programs/vdirsyncer.nix",
						},
						OptionFrom: "Home Manager",
					},
				},
			}
		}

		model, err := tui.NewTUI(getOptionsFunc)
		assert.NoError(t, err)
		tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(300, 100))

		teatest.WaitFor(
			t, tm.Output(),
			func(bts []byte) bool {
				return bytes.Contains(bts, []byte("accounts.calendar.accounts.<name>.vdirsyncer.enable"))
			},
			teatest.WithCheckInterval(time.Millisecond*100),
			teatest.WithDuration(time.Second*3),
		)

		tm.Send(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune("t"),
		})

		teatest.WaitFor(
			t, tm.Output(),
			func(bts []byte) bool {
				return bytes.Contains(bts, []byte("Description"))
			},
			teatest.WithCheckInterval(time.Millisecond*100),
			teatest.WithDuration(time.Second*3),
		)

		tm.Send(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune("q"),
		})
	})
}
