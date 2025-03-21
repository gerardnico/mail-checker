let
  # export NIX_PATH="nixpkgs=channel:nixos-24.11"
  pkgs = import <nixpkgs> { config = {}; overlays = []; };
  nodejs = pkgs.nodejs_20;
  go = pkgs.go_1_23;
  # https://semantic-release.gitbook.io/semantic-release/support/git-version
  gitMinVersion = "2.7.1";
in

assert pkgs.lib.versionAtLeast pkgs.git.version gitMinVersion || throw "Git version ${pkgs.git.version} is less than required minimum ${gitMinVersion} for semantic release";

pkgs.mkShellNoCC {
packages = with pkgs; [
  git
  go
];
shellHook = ''
    mkdir -p $out/bin
    export COREPACK_ENABLE_DOWNLOAD_PROMPT=0 && corepack enable --install-directory=$out/bin
    export PATH="$PATH:$out/bin"
    echo "Hallo from Nix"
    echo "Git Version : $(git --version)"
    echo "Go Version : $(go version)"
    export NODE_BIN="${pkgs.nodejs}/bin/node"
    echo "Node version: $(node --version)"
    echo "NPM version: $(npm --version)"
    echo "Yarn version: $(yarn --version)"
    '';
}


