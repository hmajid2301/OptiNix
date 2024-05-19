package cmd

import (
	"context"
	"database/sql"
	"log"

	// used to connect to sqlite
	_ "modernc.org/sqlite"

	"gitlab.com/hmajid2301/optinix/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func Execute(ctx context.Context, db *sql.DB) error {
	rootCmd := &cobra.Command{
		Version: "v0.1.0",
		Use:     "optinix",
		Short:   "optinix - a CLI tool to search nix options",
		Long:    `OptiNix is tool you can use on the command line to search options for NixOS, home-manager and Darwin.`,
		Example: "optinix hyprland",
		RunE: func(cmd *cobra.Command, _ []string) error {
			forceRefresh := cmd.Flags().Bool("force-refresh", false, "If set will force a refresh of the options")
			//nolint: gomnd
			limit := cmd.Flags().Int64("limit", 10, "Limit the number of results returned")

			flags := tui.Flags{
				ForceRefresh: *forceRefresh,
				Limit:        *limit,
			}
			p := tea.NewProgram(tui.New(ctx, db, flags))
			if _, err := p.Run(); err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	return rootCmd.ExecuteContext(ctx)
}
