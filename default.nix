{
  lib,
  buildGoModule,
  ...
}:
buildGoModule {
  pname = "goclacker";
  version = "1.4.3";
  src = lib.fileset.toSource {
    root = ./.;
    fileset = lib.fileset.gitTracked ./.;
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
