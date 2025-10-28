{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  ),
  mkGoEnv ? pkgs.mkGoEnv,
  gomod2nix ? pkgs.gomod2nix,
  pre-commit-hooks,
  devShellPackages ? [ ],
  ...
}:
let
  goEnv = mkGoEnv { pwd = ./.; };
  pre-commit-check = pre-commit-hooks.lib.${pkgs.system}.run {
    src = ./.;
    hooks = {
      golangci-lint.enable = true;
      gotest.enable = true;
    };
  };
in
pkgs.mkShell {
  inherit (pre-commit-check) shellHook;
  hardeningDisable = [ "all" ];
  packages =
    with pkgs;
    [
      goEnv
      gomod2nix
      docker
    ]
    ++ devShellPackages;
}
