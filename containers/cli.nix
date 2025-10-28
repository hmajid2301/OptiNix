{
  pkgs,
  package,
}:
pkgs.dockerTools.buildImage {
  name = "optinix";
  tag = "latest";
  created = "now";
  copyToRoot = pkgs.buildEnv {
    name = "image-root";
    paths = [
      package
    ];
    pathsToLink = [ "/bin" ];
  };
  config = {
    Env = [
      "OPTINIX_DB_FOLDER=/"
      "USER=root"
    ];
    Cmd = [ "${package}/bin/optinix" ];
  };
}
