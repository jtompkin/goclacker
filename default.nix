with import (fetchTarball
  "https://github.com/NixOS/nixpkgs/archive/2bfc080955153be0be56724be6fa5477b4eefabb.tar.gz") {};
  buildGoModule rec {
    pname = "goclacker";
    version = "1.4.2";
    src = fetchFromGitHub {
      owner = "jtompkin";
      repo = "goclacker";
      rev = "v${version}";
      hash = "sha256-3jELTnPFDpB5vJ+mCTV6drYIvhiPIhmQmsn2MlvaNz0=";
    };
    vendorHash = "sha256-rELkSYwqfMFX++w6e7/7suzPaB91GhbqFsLaYCeeIm4=";
  }
