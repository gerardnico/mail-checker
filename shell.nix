# default.nix, the main nix script
let
  # <nixpkgs> refers to the nixpkgs channel `channel:nixos-24.11`
  pkgs = import <nixpkgs> { config = {}; overlays = []; };
  # 2.9.0 for hooks https://stackoverflow.com/questions/39332407/git-hooks-applying-git-config-core-hookspath
  gitMinVersion = "2.9.0";
in

assert pkgs.lib.versionAtLeast pkgs.git.version gitMinVersion || throw "Git version ${pkgs.git.version} is less than required minimum ${gitMinVersion} for git hooks";

pkgs.mkShellNoCC {
      packages = with pkgs; [
        git
        go_1_23
      ];
      # Take the package’s dependencies into the environment with inputsFrom
      shellHook = ''
          echo "Hallo from Nix. Build Tools:"
          echo "Git Version        : $(git --version)"
          echo "Go Version         : $(go version)"
          '';

     }
