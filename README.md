
![Logo](docs/logo.png)

# OptiNix


A CLI tool for searching options in Nix, written in Go.

## Demo

![Demo](docs/demo.gif)

## Install

There are several ways you can install OptiNix.

### Nix flakes

Add optinix as an input to your flake, in your `flake.nix` file

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

### Nix (Coming Soon!!!)
You can install this package from nixpkgs.

```bash
nix-shell -p optinix

optinix -v
```


## Usage

```bash
optinix hyprland
```

### Nix with flakes

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
