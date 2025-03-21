# Packaging of mail-checker
# See https://nix.dev/tutorials/packaging-existing-software
# This is a function with 2 parameters
{
    buildGoModule,
    lib,
    go
}:

# manual: https://nixos.org/manual/nixpkgs/stable/#sec-language-go
# code: https://github.com/NixOS/nixpkgs/blob/master/pkgs/build-support/go/module.nix
# wiki: https://nixos.wiki/wiki/Go
buildGoModule rec {
    pname = "mail-checker";
    version = "0.0.1";
    src = ./.;
    # Directory to the `go.mod` and `go.sum` relative to the `src`
    modRoot = "./";
    # If`vendorHash` is `null`, no dependencies are fetched and the build relies on the vendor folder within the source.
    # ie go mod vendor should be run first
    vendorHash = null;
    proxyVendor = true;
    buildInputs = [ go ];
    meta = {
          description = "xxxxxx";
          homepage = "https://github.com/xxx";
          license = lib.licenses.mit;
          #maintainers = with lib.maintainers; [ kalbasit ];
    };
}
#stdenv.mkDerivation {
#    pname = "mail-checker";
#    version = "0.0.1";
#    buildInputs = [ go ];
#    pwd = ./.;
#    src = ./.;
#}