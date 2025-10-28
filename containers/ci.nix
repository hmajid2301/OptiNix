{
  pkgs,
  myPackages,
  mkGoEnv,
  gomod2nix,
  ...
}:
let
  goEnv = mkGoEnv { pwd = ../.; };
in
pkgs.dockerTools.buildImage {
  name = "optinix-dev";
  tag = "latest";
  copyToRoot = pkgs.buildEnv {
    name = "root";
    pathsToLink = [ "/bin" ];
      paths =
      with pkgs;
      [
        coreutils
        gnugrep
        bash
        curl
        git
        goEnv
        gomod2nix

        docker
        procps
        jq
        which
        tailscale
      ]
      ++ myPackages;
  };
  config = {
    Env = [
      "NIX_PAGER=cat"
      "USER=nobody"
      "SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt"
      "SSL_CERT_DIR=${pkgs.cacert}/etc/ssl/certs/"
      "CGO_ENABLED=0"
    ];
  };
}
