{
  pkgs,
  package,
}: let
  nonRootShadowSetup = {
    user,
    uid,
    gid ? uid,
  }:
    with pkgs; [
      (
        writeTextDir "etc/shadow" ''
          root:!x:::::::
          ${user}:!:::::::
        ''
      )
      (
        writeTextDir "etc/passwd" ''
          root:x:0:0::/root:${runtimeShell}
          ${user}:x:${toString uid}:${toString gid}::/home/${user}:
        ''
      )
      (
        writeTextDir "etc/group" ''
          root:x:0:
          ${user}:x:${toString gid}:
        ''
      )
      (
        writeTextDir "etc/gshadow" ''
          root:x::
          ${user}:x::
        ''
      )
    ];
in
  pkgs.dockerTools.buildImage {
    name = "optinix";
    tag = "0.1";
    created = "now";
    extraCommands = ''
      mkdir -p /tmp
    '';
    copyToRoot = pkgs.buildEnv {
      name = "image-root";
      paths =
        [
          pkgs.nix
          pkgs.bash
          pkgs.coreutils
          package
        ]
        ++ nonRootShadowSetup {
          uid = 1000;
          user = "nixbld";
        };
      pathsToLink = ["/bin" "/etc" "/tmp"];
    };
    config = {
      Env = [
        "OPTINIX_DB_FOLDER=/"
        "USER=root"
      ];
      Cmd = ["${package}/bin/optinix"];
    };
  }
