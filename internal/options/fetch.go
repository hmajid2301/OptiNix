package options

import (
	"context"
	"fmt"
	"path/filepath"
)

type CmdExecutor interface {
	// TODO: rename
	// TODO: Add context
	Executor(expression string) (string, error)
}

type Reader interface {
	// TODO: Add context
	Read(r string) ([]byte, error)
}

type Sources struct {
	NixOS       string
	HomeManager string
	Darwin      string
}

type Fetcher struct {
	executor CmdExecutor
	reader   Reader
}

func NewFetcher(executor CmdExecutor, reader Reader) Fetcher {
	return Fetcher{executor: executor, reader: reader}
}

func (f Fetcher) Fetch(ctx context.Context, sources Sources) ([]Option, error) {
	fmt.Println("Fetching options", ctx)
	var options []Option
	for _, source := range []string{sources.NixOS, sources.HomeManager, sources.Darwin} {
		var path string
		var err error

		switch source {
		case sources.NixOS:
			path, err = f.GetNixosDocPath(source)
		case sources.HomeManager:
			path, err = f.GetHMDocPath(source)
		case sources.Darwin:
			path, err = f.GetDarwinDocPath(source)
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

// TODO: remind user to add hm to nix-channel
// nix-channel --add https://github.com/nix-community/home-manager/archive/master.tar.gz home-manager
// nix-channel --update
func (f Fetcher) GetHMDocPath(expression string) (string, error) {
	output, err := f.executor.Executor(expression)
	if err != nil {
		return "", err
	}

	path := filepath.Join(output, "share/doc/home-manager/options.json")
	return path, nil
}

func (f Fetcher) GetNixosDocPath(expression string) (string, error) {
	output, err := f.executor.Executor(expression)
	return output, err
}

func (f Fetcher) GetDarwinDocPath(expression string) (string, error) {
	output, err := f.executor.Executor(expression)
	return output, err
}
