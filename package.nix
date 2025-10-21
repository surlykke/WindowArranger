{
  lib,
  buildGoModule,
  fetchFromGitHub,
}:

buildGoModule rec {
  pname = "windowarranger";
  version = "0.1";
  src = lib.cleanSource ./.;
	vendorHash = null;
  meta = {
    description = "Window arranger for sway";
    homepage = "https://github.com/surlykke/windowarranger";
    license = lib.licenses.gpl2;
    maintainers = with lib.maintainers; [ surlykke ];
  };
}
