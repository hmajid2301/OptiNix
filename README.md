
![Logo](docs/logo.png)

# OptiNix

A CLI tool for searching options in Nix, written in Go.

## Demo

![Demo](docs/demo.gif)

## Install

There are several ways you can install OptiNix.

### Nix flakes

Add OptiNix as an input to your flake, in your `flake.nix` file

```nix
{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    optinix.url = "gitlab:hmajid2301/optinix";
  };
  outputs = {};
}
```

Then you can install the package in your nix config like in a `devShell`

```nix
{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    optinix.url = "gitlab:hmajid2301/optinix";
  };
  outputs = {
    self,
    nixpkgs,
    flake-utils,
    optinix,
    ...
  }: (
    flake-utils.lib.eachDefaultSystem
    (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          optinix.packages.${system}.default
        ];
      };
    })
  );
}
```


### Go

You can install using `go`

```bash
go install gitlab.com/hmajid2301/optinix
```

<!-- ### Nix (Coming Soon.) -->
<!-- You can install this package from nixpkgs. -->
<!---->
<!-- ```bash -->
<!-- nix-shell -p optinix -->
<!---->
<!-- optinix -v -->
<!-- ``` -->


## Usage

```bash
optinix --help
OptiNix is tool you can use on the command line to search options for NixOS, home-manager and Darwin.

Usage:
  optinix [flags]

Examples:
optinix hyprland

Flags:
      --force-refresh   If set will force a refresh of the options
  -h, --help            help for optinix
      --limit int       Limit the number of results returned (default 10)
      --no-tui          If set will not show the TUI and just print the options to stdout
  -v, --version         version for optinix


# Example in the demo above
optinix hyprland
```

### Without TUI

By default, the tool will render a list using bubble tea, if you want just plain text output you can add the `--no-tui` flag.

```bash
optinix podman --no-tui


Option: 0
Name: virtualisation.podman.networkSocket.server
Type: value "ghostunnel" (singular enum)
Default:
Example: "ghostunnel"
From: NixOS
Sources: [/nix/store/rhg90jpryc286xn9xjy6qjiaap6pjgdc-source/nixos/modules/virtualisation/podman/network-socket-ghostunnel.nix /nix/store/rhg90jpryc286xn9xjy6qjiaap6pjgdc-source/nixos/modules/virtualisation/podman/network-socket.nix]


Option: 1
Name: virtualisation.podman.networkSocket.enable
Type: boolean
Default: false
Example:
From: NixOS
Sources: [/nix/store/rhg90jpryc286xn9xjy6qjiaap6pjgdc-source/nixos/modules/virtualisation/podman/network-socket.nix]


Option: 2
Name: virtualisation.podman.enable
Type: boolean
Default: false
Example:
From: NixOS
Sources: [/nix/store/rhg90jpryc286xn9xjy6qjiaap6pjgdc-source/nixos/modules/virtualisation/podman/default.nix]

# ...
```

### With FZF
You can integrate this tool with FZF as well, like so:

```bash
optinix --no-tui | rg "Name: " | cut -d' ' -f2 | fzf --preview="optinix --no-tui '{}'"
```

### Nix with flakes

If you are running Nix and have flakes enabled in your configuration, you can run the tool like the command below. Without needing to "install" OptiNix.

```bash
nix run 'gitlab.com/hmajid2301/optinix' hyprland
```

## Supported Sources

- NixOS
- Home Manager
- Darwin

## Inspired By
- Manix: https://github.com/nix-community/manix
  - When I started this project back last year, this project was not working but has since been fixed and had new features added
