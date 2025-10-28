{
  description = "Development environment for OptiNix";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    pre-commit-hooks.url = "github:cachix/pre-commit-hooks.nix";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      gomod2nix,
      pre-commit-hooks,
      ...
    }:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        myPackages = with pkgs; [
          go_1_25

          golangci-lint
          gotools
          gotestsum
          gocover-cobertura
          go-task

          sqlc
          sqlfluff
        ];

        devShellPackages =
          with pkgs;
          myPackages
          ++ [
            go-junit-report
            goreleaser
            vhs
            gum
          ];
      in
      rec {
        packages.default = pkgs.callPackage ./. {
          inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
        };
        devShells.default = pkgs.callPackage ./shell.nix {
          inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
          inherit pre-commit-hooks;
          inherit devShellPackages;
        };
        packages.container = pkgs.callPackage ./containers/service.nix {
          package = packages.default;
        };
        packages.container-ci = pkgs.callPackage ./containers/ci.nix {
          inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
          inherit myPackages;
        };
      }
    ));
}
