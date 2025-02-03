{
  description = "Wait4X allows you to wait for a port or a service to enter the requested state.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
    nixpkgs-unstable.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    nixpkgs-unstable,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
      unstable = nixpkgs-unstable.legacyPackages.${system};
      packageName = "wait4x";
    in {
      formatter = pkgs.alejandra;
      devShells.default = pkgs.mkShell {
        name = packageName;
        buildInputs = with pkgs; [
          go
        ];
      };
      packages.default = pkgs.buildGoModule {
        pname = packageName;
        version = "${self.shortRev or self.dirtyShortRev or "dirty"}";
        src = self;
        vendorHash = "sha256-KtEOLLsbTfgaXy/0aj5zT5qbgW6qBFMuU3EnnXRu+Ig=";
        doCheck = false;
      };
    });
}
