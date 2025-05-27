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
buildGoModule rec {
  pname = "goclacker";
  version = "1.4.2";
  src = fs.toSource {
    root = ./.;
    fileset = src;
  };
  vendorHash = lib.fakeHash;
}
