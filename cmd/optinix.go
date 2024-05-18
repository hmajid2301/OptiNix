package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"

	// used to connect to sqlite
	_ "modernc.org/sqlite"

	"gitlab.com/hmajid2301/optinix/internal/options"
	"gitlab.com/hmajid2301/optinix/internal/options/store"

	"github.com/spf13/cobra"
)

type Flags struct {
	forceRefresh bool
}

func Execute(ctx context.Context, db *sql.DB) error {
	rootCmd := &cobra.Command{
		Version: "v0.1.0",
		Use:     "optinix",
		Short:   "optinix - a CLI tool to search nix options",
		Long:    `OptiNix is tool you can use on the command line to search options for NixOS, home-manager and Darwin.`,
		Args:    cobra.ExactArgs(1),
		Example: "optinix hyprland",
		RunE: func(cmd *cobra.Command, args []string) error {
			forceRefresh := cmd.Flags().Bool("force-refresh", false, "If set will force a refresh of the options")

			// p := tea.NewProgram(tui.New(ctx, db))
			// if _, err := p.Run(); err != nil {
			// 	log.Fatal(err)
			// }
			flags := Flags{
				forceRefresh: *forceRefresh,
			}
			err := FindOptions(ctx, db, args[0], flags)
			return err
		},
	}

	return rootCmd.ExecuteContext(context.Background())
}

func FindOptions(ctx context.Context, db *sql.DB, optionName string, flags Flags) (err error) {
	s, err := store.NewStore(db)
	if err != nil {
		return err
	}

	// TODO: should this be setup with constructors
	cmdExecutor := NixCmdExecutor{}
	nixReader := NixReader{}
	fetcher := options.NewFetcher(cmdExecutor, nixReader)

	opt := options.NewOptions(s, fetcher)

	// TODO: should I read file and evalute expression?
	nixosPath, err := nixReader.Read("nix/nixos-options.nix")
	if err != nil {
		return err
	}

	homeManagerPath, err := nixReader.Read("nix/hm-options.nix")
	if err != nil {
		return err
	}

	darwinPath, err := nixReader.Read("nix/darwin-options.nix")
	if err != nil {
		return err
	}

	sources := options.Sources{
		NixOS:       string(nixosPath),
		HomeManager: string(homeManagerPath),
		Darwin:      string(darwinPath),
	}
	err = opt.SaveOptions(ctx, sources, flags.forceRefresh)
	if err != nil {
		return err
	}

	opts, err := opt.GetOptions(ctx, optionName)
	if err != nil {
		return err
	}

	// TODO: format this nicely
	for _, o := range opts {
		fmt.Println(o.Name)
		fmt.Println(o.Type)
		fmt.Println(o.Description)
		fmt.Println(o.DefaultValue)
		fmt.Println(o.Example)

		for _, s := range o.Sources {
			fmt.Println(s)
		}
	}

	return nil
}

type NixReader struct{}

func (f NixReader) Read(pathToExpression string) ([]byte, error) {
	nixExpression, err := os.ReadFile(pathToExpression)
	return nixExpression, err
}

type NixCmdExecutor struct{}

func (n NixCmdExecutor) Executor(expression string) (string, error) {
	cmd := exec.Command("nix-build", "-E", expression)
	cmd.Env = append(cmd.Env,
		"NIXPKGS_ALLOW_UNFREE=1",
		"NIXPKGS_ALLOW_BROKEN=1",
		"NIXPKGS_ALLOW_INSECURE=1",
		"NIXPKGS_ALLOW_UNSUPPORTED_SYSTEM=1",
	)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	trimmedOuput := strings.TrimSpace(string(output))
	return trimmedOuput, nil
}
