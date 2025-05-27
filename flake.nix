{
  description = "Command line reverse Polish notation calculator";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      inherit (nixpkgs) lib;
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = lib.genAttrs systems;
    in
    {
      packages = forAllSystems (system: rec {
        goclacker = nixpkgs.legacyPackages.${system}.callPackage ./. { };
        default = goclacker;
      });
      apps = forAllSystems (system: rec {
        goclacker = {
          type = "app";
          program = "${self.packages.${system}.goclacker}/bin/goclacker";
        };
        default = goclacker;
      });
      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        rec {
          goclacker = pkgs.mkShell { packages = [ pkgs.go ]; };
          default = goclacker;
        }
      );
    };
}
