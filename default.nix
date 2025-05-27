{ lib, buildGoModule, ... }:
let
  fs = lib.fileset;
  src = fs.difference (fs.gitTracked ./.) (
    fs.unions [
      ./.envrc
      ./flake.lock
      (fs.fileFilter (file: file.hasExt "md") ./.)
      (fs.fileFilter (file: file.hasExt "nix") ./.)
    ]
  );
in
buildGoModule {
  pname = "goclacker";
  version = "1.4.2";
  src = fs.toSource {
    root = ./.;
    fileset = src;
  };
  vendorHash = "sha256-rELkSYwqfMFX++w6e7/7suzPaB91GhbqFsLaYCeeIm4=";
}
