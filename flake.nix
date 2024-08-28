{
  description = "Wait for a port or a service to enter the requested state.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }@inputs: flake-utils.lib.eachDefaultSystem (system: {
    packages = let
      inherit (nixpkgs) lib;
      inherit (nixpkgs.legacyPackages.${system}) buildGoModule;
    in {
      default = buildGoModule {
        pname = "wait4x";
        version = builtins.substring 0 8 self.lastModifiedDate;

        src = self;
        # don't know why but nix is unhappy about vendorHash = null
        vendorHash = "sha256-Jp2IUvkcqLcJk0a5A79SQTjqAkmIEVc9Ove3rMkkWuI=";

        # Nix doesn't allow network access during tests, so belive in the existing tests
        doCheck = false;
      };
    };
  });
}
