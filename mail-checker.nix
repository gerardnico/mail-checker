# Packaging of mail-checker
# See https://nix.dev/tutorials/packaging-existing-software
# This is a function with 2 parameters
{
    stdenv
}:

stdenv.mkDerivation {
    pname = "mail-checker";
    version = "0.0.1";
    pwd = ./.;
    src = ./.;
}