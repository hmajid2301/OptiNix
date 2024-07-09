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
  buildGoApplication ? pkgs.buildGoApplication,
}:
buildGoApplication {
  pname = "optinix";
  version = "0.1";
  go = pkgs.go_1_22;
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  # postInstall = ''
  #   mkdir -p $out/share/bash-completion/completions
  #   mkdir -p $out/share/zsh/site-functions
  #   mkdir -p $out/share/fish/vendor_completions.d
  #   $out/bin/$pname completion bash > $out/share/bash-completion/completions/$pname
  #   $out/bin/$pname completion zsh > $out/share/zsh/site-functions/_$pname
  #   $out/bin/$pname completion fish > $out/share/fish/vendor_completions.d/$pname.fish
  # '';
}
