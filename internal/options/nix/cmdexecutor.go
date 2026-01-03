package nix

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type CmdExecutor struct{}

func NewCmdExecutor() CmdExecutor {
	return CmdExecutor{}
}

func (c CmdExecutor) ExecuteCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = os.Environ()
	return cmd
}

func (CmdExecutor) Execute(ctx context.Context, expression string) (string, error) {
	cmd := exec.CommandContext(ctx, "nix-build", "--no-out-link", "-E", expression)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		"NIXPKGS_ALLOW_UNFREE=1",
		"NIXPKGS_ALLOW_BROKEN=1",
		"NIXPKGS_ALLOW_INSECURE=1",
		"NIXPKGS_ALLOW_UNSUPPORTED_SYSTEM=1",
	)

	if nixPath, ok := os.LookupEnv("NIX_PATH"); ok {
		cmd.Env = append(cmd.Env, fmt.Sprintf("NIX_PATH=%s", nixPath))
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}

	trimmedOuput := strings.TrimSpace(string(output))
	return trimmedOuput, nil
}
