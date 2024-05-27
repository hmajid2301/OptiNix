package nix

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type CmdExecutor struct{}

func NewCmdExecutor() CmdExecutor {
	return CmdExecutor{}
}

func (CmdExecutor) Execute(ctx context.Context, expression string) (string, error) {
	cmd := exec.CommandContext(ctx, "nix-build", expression)
	cmd.Env = append(cmd.Env,
		"NIXPKGS_ALLOW_UNFREE=1",
		"NIXPKGS_ALLOW_BROKEN=1",
		"NIXPKGS_ALLOW_INSECURE=1",
		"NIXPKGS_ALLOW_UNSUPPORTED_SYSTEM=1",
		// TODO: why does this needed on NixOS but not on Ubuntu
		"NIX_PATH=/etc/nix/inputs",
		"--no-out-link",
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}

	trimmedOuput := strings.TrimSpace(string(output))
	return trimmedOuput, nil
}
