{
  description = "Command line reverse Polish notation calculator";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.05";
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
