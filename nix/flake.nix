{
  description = "OptiNix internal flake for fetching NixOS/Home Manager/Darwin options";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    darwin = {
      url = "github:LnL7/nix-darwin";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = {
    self,
    nixpkgs,
    home-manager,
    darwin,
  }: let
    systems = ["x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin"];
    forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f system);
  in {
    packages = forAllSystems (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      # NixOS options
      nixos-options = let
        eval = import (nixpkgs + "/nixos/lib/eval-config.nix") {
          inherit system;
          modules = [];
        };
        opts = (pkgs.nixosOptionsDoc {options = eval.options;}).optionsJSON;
      in
        pkgs.runCommand "nixos-options.json" {inherit opts;}
        "cp $opts/share/doc/nixos/options.json $out";

      # Home Manager options
      home-manager-options = let
        lib = import (home-manager + "/modules/lib/stdlib-extended.nix") pkgs.lib;
        hmargs = {
          release = "unstable";
          isReleaseBranch = false;
          inherit pkgs lib;
        };
        docs = import (home-manager + "/docs") hmargs;
        opts = (
          if builtins.isFunction docs
          then docs hmargs
          else docs
        )
        .options
        .json;
      in
        pkgs.runCommand "home-manager-options.json" {inherit opts;}
        "cp $opts/share/doc/home-manager/options.json $out";

      # Darwin options
      darwin-options = let
        eval = import darwin {
          inherit system pkgs;
          configuration = {...}: {
            system.stateVersion = 5;
          };
        };
        opts = eval.config.system.build.manual.optionsJSON;
      in
        pkgs.runCommand "darwin-options.json" {inherit opts;}
        "cp $opts/share/doc/darwin/options.json $out";
    });
  };
}
