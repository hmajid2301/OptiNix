{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
, mkGoEnv ? pkgs.mkGoEnv
, gomod2nix ? pkgs.gomod2nix
}:

let
  goEnv = mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  hardeningDisable = [ "all" ];
  packages = [
    goEnv
    gomod2nix
    pkgs.gopls
    pkgs.golangci-lint
    pkgs.go_1_21
    pkgs.gotools
    pkgs.go-junit-report
    pkgs.go-task
    pkgs.delve
  ];
}
