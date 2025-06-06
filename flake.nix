{
  description = "Command line reverse Polish notation calculator";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.05";
  };

  outputs =
    { self, nixpkgs, ... }:
    let
      inherit (nixpkgs) lib;
      allSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems =
        f:
        lib.genAttrs allSystems (
          system:
          f {
            inherit system;
            pkgs = import nixpkgs { inherit system; };
          }
        );
    in
    {
      packages = forAllSystems (
        { pkgs, system }:
        {
          goclacker = pkgs.callPackage self { };
          default = self.packages.${system}.goclacker;
        }
      );
      apps = forAllSystems (
        { system, ... }:
        {
          goclacker = {
            type = "app";
            program = lib.getExe self.packages.${system}.goclacker;
            meta.description = "Command line reverse Polish notation calculator";
          };
          default = self.apps.${system}.goclacker;
        }
      );
      devShells = forAllSystems (
        { pkgs, system }:
        {
          goclacker = pkgs.mkShell { packages = [ pkgs.go ]; };
          default = self.devShells.${system}.goclacker;
        }
      );
      overlays = {
        goclacker = final: prev: {
          goclacker = final.callPackage self { };
        };
        default = self.overlays.goclacker;
      };
      homeModules = {
        goclacker =
          { pkgs, ... }:
          let
            inherit (pkgs.stdenv.hostPlatform) system;
          in
          {
            imports = [ ./modules/home-module.nix ];
            programs.goclacker.package = self.packages.${system}.goclacker;
          };
        default = self.homeModules.goclacker;
      };
    };
}
