# default.nix
let
  pkgs = import <nixpkgs> { config = {}; overlays = []; };
in
{
  mail-checker = pkgs.callPackage ./mail-checker.nix { };
}