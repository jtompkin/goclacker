{
  lib,
  buildGoModule,
  fetchFromGitHub,
  ...
}:
buildGoModule rec {
  pname = "goclacker";
  version = "1.4.3";
  src = fetchFromGitHub {
    owner = "jtompkin";
    repo = pname;
    tag = "v${version}";
    hash = lib.fakeHash;
  };
  vendorHash = "sha256-rELkSYwqfMFX++w6e7/7suzPaB91GhbqFsLaYCeeIm4=";
  meta = {
    description = "Command line reverse Polish notation calculator";
    homepage = "https://github.com/jtompkin/goclacker";
    license = lib.licenses.mit;
    platforms = lib.platforms.all;
    mainProgram = "goclacker";
  };
}
