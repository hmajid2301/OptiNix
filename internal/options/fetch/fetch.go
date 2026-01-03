package fetch

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
)

type Exectutor interface {
	Execute(ctx context.Context, path string) (string, error)
}

type Reader interface {
	Read(path string) ([]byte, error)
}

type Fetcher struct {
	nixCmdExecutor Exectutor
	reader         Reader
	messenger      entities.Messager
	logger         *slog.Logger
}

func NewFetcher(executor Exectutor, reader Reader, messager entities.Messager, logger *slog.Logger) Fetcher {
	return Fetcher{nixCmdExecutor: executor, reader: reader, messenger: messager, logger: logger}
}

func (f Fetcher) Fetch(ctx context.Context, sources entities.Sources) ([]entities.Option, error) {
	var options []entities.Option

	// Build options from internal flake
	type optionSource struct {
		name        string
		packageName string
		optionFrom  string
	}

	optionSources := []optionSource{
		{"NixOS", "nixos-options", "NixOS"},
		{"Home Manager", "home-manager-options", "Home Manager"},
		{"Darwin", "darwin-options", "Darwin"},
	}

	totalSources := len(optionSources)
	for idx, src := range optionSources {
		f.messenger.Send(fmt.Sprintf("Fetching %s options... (%d/%d)", src.name, idx+1, totalSources))

		nixExpr := fmt.Sprintf(`(builtins.getFlake (toString ./nix)).packages.%s.%s`,
			"${builtins.currentSystem}", src.packageName)

		path, err := f.nixCmdExecutor.Execute(ctx, nixExpr)
		if err != nil {
			f.messenger.Send(fmt.Sprintf("failed to get %s options from internal flake", src.name))
			return nil, err
		}

		f.logger.Debug("Got path", "path", path)

		contents, err := f.reader.Read(strings.TrimSpace(path))
		if err != nil {
			return nil, err
		}

		f.messenger.Send(fmt.Sprintf("Parsing %s options... (%d/%d)", src.name, idx+1, totalSources))
		opts, err := ParseOptions(contents, src.optionFrom)
		if err != nil {
			return nil, err
		}
		options = append(options, opts...)
	}

	return options, nil
}
