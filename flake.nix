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
      packages = forAllSystems (system: {
        goclacker = nixpkgs.legacyPackages.${system}.callPackage ./. { };
        default = self.packages.${system}.goclacker;
      });
      apps = forAllSystems (system: {
        goclacker = {
          type = "app";
          program = lib.getExe self.packages.${system}.goclacker;
        };
        default = self.apps.${system}.goclacker;
      });
      overlays = {
        goclacker = final: prev: {
          goclacker = final.callPackage ./. { };
        };
        default = self.overlays.goclacker;
      };
      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          goclacker = pkgs.mkShell { packages = [ pkgs.go ]; };
          default = self.devShells.${system}.goclacker;
        }
      );
    };
}
