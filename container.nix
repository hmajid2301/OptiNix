{
  pkgs,
  package,
}: let
  nixFromDockerHub = pkgs.dockerTools.pullImage {
    imageName = "nixos/nix";
    imageDigest = "sha256:4999c663e350df27f0c77d4636c324f7093ff0bfce0a123aa249b8c74289c78b";
    sha256 = "sha256-o8k+2PSCfjkvR9keLG71fUlMSti+z5+09Vj+bEadscM=";
    finalImageTag = "2.22.2";
    finalImageName = "nix";
  };
in
  pkgs.dockerTools.buildImage {
    name = "optinix";
    tag = "latest";
    created = "now";
    fromImage = nixFromDockerHub;
    copyToRoot = pkgs.buildEnv {
      name = "image-root";
      paths = [
        package
      ];
      pathsToLink = ["/bin"];
    };
    diskSize = 2048;
    buildVMMemorySize = 2048;

    # runAsRoot = ''
    #   mkdir /tmp
    # '';

    config = {
      Env = [
        "OPTINIX_DB_FOLDER=/"
        "USER=root"
      ];
      Cmd = ["${package}/bin/optinix"];
    };
  }
