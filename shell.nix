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
, pre-commit-hooks
, ...
}:

let
  goEnv = mkGoEnv { pwd = ./.; };
  pre-commit-check = pre-commit-hooks.lib.${pkgs.system}.run {
    src = ./.;
    hooks = {
      golangci-lint = {
        enable = true;
        name = "golangci-lint";
        description = "Lint my golang code";
        files = "\.go$";
        entry = "task lint -- --new-from-rev HEAD --fix";
        require_serial = true;
        pass_filenames = false;
      };
    };
  };
in

pkgs.mkShell {
  inherit (pre-commit-check) shellHook;
  hardeningDisable = [ "all" ];
  packages = [
    goEnv
    gomod2nix
    pkgs.golangci-lint
    pkgs.go_1_21
    pkgs.gotools
    pkgs.go-junit-report
    pkgs.go-task
    pkgs.delve
  ];
}
