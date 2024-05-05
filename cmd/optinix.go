package cmd

import (
	"context"
	"database/sql"
	"log"

	// used to connect to sqlite
	_ "modernc.org/sqlite"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"gitlab.com/hmajid2301/optinix/internal/tui"
)

func Execute(ctx context.Context, db *sql.DB) error {
	rootCmd := &cobra.Command{
		// Version: "v0.1.0",
		Use:   "optinix",
		Short: "optinix - a CLI tool to search nix options",
		Long: `OptiNix is tool you can use on the command line to search options for both NixOS and home-manager
	rather than needing to go to a website i.e. nixos.org or mynixos.com.`,
		// Args:    cobra.ExactArgs(1),
		Example: "optinix",
		RunE: func(_ *cobra.Command, _ []string) error {
			p := tea.NewProgram(tui.New(ctx, db))
			if _, err := p.Run(); err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	return rootCmd.ExecuteContext(context.Background())
}
