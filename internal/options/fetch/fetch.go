package fetch

import (
	"context"
	"fmt"
	"path/filepath"

	"gitlab.com/hmajid2301/optinix/internal/options/entities"
)

type Exectutor interface {
	Execute(ctx context.Context, path string) (string, error)
}

type Reader interface {
	Read(path string) ([]byte, error)
}

type Fetcher struct {
	NixCmdExecutor Exectutor
	reader         Reader
}

func NewFetcher(executor Exectutor, reader Reader) Fetcher {
	return Fetcher{NixCmdExecutor: executor, reader: reader}
}

func (f Fetcher) Fetch(ctx context.Context, sources entities.Sources) ([]entities.Option, error) {
	var options []entities.Option
	for _, source := range []string{sources.NixOS, sources.HomeManager, sources.Darwin} {
		var path string
		var err error

		switch source {
		case sources.NixOS:
			path, err = f.GetNixosDocPath(ctx, source)
		case sources.HomeManager:
			path, err = f.GetHMDocPath(ctx, source)
			if err != nil {
				err = fmt.Errorf(`failed to get home-manager options, try to run:\n`+
					`nix-channel --add https://github.com/nix-community/home-manager/archive/master.tar.gz home-manager\n`+
					`nix-channel --update\n\n`+
					`%s`, err)
			}
		case sources.Darwin:
			path, err = f.GetDarwinDocPath(ctx, source)
		}

		if err != nil {
			return nil, err
		}

		contents, err := f.reader.Read(path)
		if err != nil {
			return nil, err
		}

		opts, err := ParseOptions(contents)
		if err != nil {
			return nil, err
		}
		options = append(options, opts...)
	}

	return options, nil
}

func (f Fetcher) GetHMDocPath(ctx context.Context, expression string) (string, error) {
	output, err := f.NixCmdExecutor.Execute(ctx, expression)
	if err != nil {
		return "", err
	}

	path := filepath.Join(output, "share/doc/home-manager/options.json")
	return path, nil
}

func (f Fetcher) GetNixosDocPath(ctx context.Context, expression string) (string, error) {
	output, err := f.NixCmdExecutor.Execute(ctx, expression)
	return output, err
}

func (f Fetcher) GetDarwinDocPath(ctx context.Context, expression string) (string, error) {
	output, err := f.NixCmdExecutor.Execute(ctx, expression)
	return output, err
}
