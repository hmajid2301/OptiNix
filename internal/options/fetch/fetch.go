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

type Messager interface {
	Send(msg string)
}

type Fetcher struct {
	nixCmdExecutor Exectutor
	reader         Reader
	Messager       Messager
}

func NewFetcher(executor Exectutor, reader Reader, messager Messager) Fetcher {
	return Fetcher{nixCmdExecutor: executor, reader: reader, Messager: messager}
}

func (f Fetcher) Fetch(ctx context.Context, sources entities.Sources) ([]entities.Option, error) {
	var options []entities.Option
	for _, source := range []string{sources.NixOS, sources.HomeManager, sources.Darwin} {
		var path string
		var err error
		var optionFrom string

		switch source {
		case sources.NixOS:
			optionFrom = "NixOS"
			f.Messager.Send("Trying to fetch NixOS options")
			path, err = f.GetNixosDocPath(ctx, source)
		case sources.HomeManager:
			optionFrom = "Home Manager"
			f.Messager.Send("Trying to fetch Home Manager options")
			path, err = f.GetHMDocPath(ctx, source)
			if err != nil {
				f.Messager.Send(`failed to get home-manager options, try to run:\n` +
					`nix-channel --add https://github.com/nix-community/home-manager/archive/master.tar.gz home-manager\n` +
					`nix-channel --update\n\n`)
			}
		case sources.Darwin:
			optionFrom = "Darwin"
			f.Messager.Send("Trying to fetch Darwin options")
			path, err = f.GetDarwinDocPath(ctx, source)
		}

		f.Messager.Send(fmt.Sprintf("err: %s", err))
		if err != nil {
			return nil, err
		}

		contents, err := f.reader.Read(path)
		if err != nil {
			return nil, err
		}

		opts, err := ParseOptions(contents, optionFrom)
		if err != nil {
			return nil, err
		}
		options = append(options, opts...)
	}

	return options, nil
}

func (f Fetcher) GetHMDocPath(ctx context.Context, expression string) (string, error) {
	output, err := f.nixCmdExecutor.Execute(ctx, expression)
	if err != nil {
		return "", err
	}

	path := filepath.Join(output, "share/doc/home-manager/options.json")
	return path, nil
}

func (f Fetcher) GetNixosDocPath(ctx context.Context, expression string) (string, error) {
	output, err := f.nixCmdExecutor.Execute(ctx, expression)
	return output, err
}

func (f Fetcher) GetDarwinDocPath(ctx context.Context, expression string) (string, error) {
	output, err := f.nixCmdExecutor.Execute(ctx, expression)
	return output, err
}
