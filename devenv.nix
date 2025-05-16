{
  pkgs,
  lib,
  config,
  inputs,
  ...
}:

let
  pkgs-unstable = import inputs.unstable { system = pkgs.stdenv.system; };
in
{
  # https://devenv.sh/basics/
  env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = [
    pkgs.git
    pkgs.nixd
  ];

  # https://devenv.sh/languages/
  languages.go = {
    enable = true;
    package = pkgs-unstable.go;
  };

  # https://devenv.sh/processes/
  # processes.cargo-watch.exec = "cargo-watch";

  # https://devenv.sh/services/
  # Test db
  services.postgres =  {
        enable = true;
        listen_addresses = "127.0.0.1";
        port = 5533;
        initialScript = ''
            CREATE ROLE postgres SUPERUSER;
        '';
        initialDatabases = [
            {
                name = "hamed";
                user = "hamed";
                pass = "1234";
            }
        ];
    };

  # https://devenv.sh/scripts/
  scripts.hello.exec = ''
    echo hello from $GREET
  '';

  # export CGO_LDFLAGS="-L${pkgs.unixODBC}/lib"
  # export CGO_CFLAGS="-I${pkgs.unixODBC}/include"
  enterShell = ''
    hello
    git --version
  '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';

  # https://devenv.sh/git-hooks/
  # git-hooks.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/
}
