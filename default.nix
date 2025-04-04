with import (fetchTarball
  "https://github.com/NixOS/nixpkgs/archive/2bfc080955153be0be56724be6fa5477b4eefabb.tar.gz") {};
  buildGoModule rec {
    pname = "goclacker";
    version = "1.4.1";
    src = fetchFromGitHub {
      owner = "jtompkin";
      repo = "goclacker";
      rev = "v${version}";
      hash = "sha256-+kJWNFlSgrgv49RT9/AHWfj1rHXgna7JphN+Be2pNpw=";
    };
    vendorHash = "sha256-rELkSYwqfMFX++w6e7/7suzPaB91GhbqFsLaYCeeIm4=";
  }
