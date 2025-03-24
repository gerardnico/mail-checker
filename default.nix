# default.nix
# The main nix script
# Tuto: https://nix.dev/guides/recipes/sharing-dependencies
let
  # <nixpkgs> refers to NIX_PATH
  # export NIX_PATH="nixpkgs=channel:nixos-24.11"
  pkgs = import <nixpkgs> { config = {}; overlays = []; };
  # 2.7.1 for semantic release https://semantic-release.gitbook.io/semantic-release/support/git-version
  # 2.9.0 for hooks https://stackoverflow.com/questions/39332407/git-hooks-applying-git-config-core-hookspath
  gitMinVersion = "2.9.0";
  mail-checker = pkgs.callPackage ./build.nix { go = pkgs.go_1_23; };
in

assert pkgs.lib.versionAtLeast pkgs.git.version gitMinVersion || throw "Git version ${pkgs.git.version} is less than required minimum ${gitMinVersion} for git hooks";

# Returned object
{

  # Create a mail-checker
  # that you can run with nix-build -A mail-checker
  inherit mail-checker;
  # Create a shell function (used in shell.nix)
  shell = pkgs.mkShellNoCC {
      packages = with pkgs; [
        git
        nodejs
        nodePackages."@commitlint/cli"
        nodePackages."@commitlint/config-conventional"
        jreleaser-cli
      ];
      # Take the package’s dependencies into the environment with inputsFrom
      inputsFrom = [ mail-checker ];
      shellHook = ''
          echo "Hallo from Nix"

          echo "Git Version        : $(git --version)"

          echo "Go Version         : $(go version)"

          echo "JReleaser version  : $(jreleaser --version | grep jreleaser | cut -d ' ' -f '2-')"

          echo "CommitLint version : $(commitlint -version)"

          # Node for commitlint
          export NODE_BIN="${pkgs.nodejs}/bin/node"
          #echo "Node version: $(node --version)"
          #echo "NPM version: $(npm --version)"
          #echo "NPX version: $(npx --version)"

          #mkdir -p $out/bin
          #export COREPACK_ENABLE_DOWNLOAD_PROMPT=0 && corepack enable --install-directory=$out/bin
          #export PATH="$PATH:$out/bin"
          #echo "Yarn version: $(yarn --version)"
          '';

     };
}