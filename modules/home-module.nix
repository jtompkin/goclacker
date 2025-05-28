{
  config,
  pkgs,
  lib,
  ...
}:
let
  cfg = config.programs.goclacker;
  inherit (lib) mkIf;
in
{
  options.programs.goclacker.enable = lib.mkEnableOption "goclacker RPN calculator";
  options.programs.goclacker.package = lib.mkPackageOption pkgs "goclacker" { };
  config = mkIf cfg.enable {
    home.packages = [ cfg.package ];
  };
}
