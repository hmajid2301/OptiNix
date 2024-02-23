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

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    gomod2nix,
    pre-commit-hooks,
    ...
  }: (
    flake-utils.lib.eachDefaultSystem
    (system: let
      pkgs = nixpkgs.legacyPackages.${system};

      # The current default sdk for macOS fails to compile go projects, so we use a newer one for now.
      # This has no effect on other platforms.
      callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;
    in rec {
      packages.default = callPackage ./. {
        inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
      };
      devShells.default = callPackage ./shell.nix {
        inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
        inherit pre-commit-hooks;
      };
      packages.container = pkgs.callPackage ./container.nix {package = packages.default;};
    })
  );
}
