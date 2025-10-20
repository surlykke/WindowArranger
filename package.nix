# First, the named parameters.
{
  # you almost always want to depend on lib
  lib,

  # these are Nixpkgs functions that your package uses
  buildGoModule,
  fetchFromGitHub,
  # more dependencies would go here...
}: # this means end of named parameters

# Now, the definition of your package.
# This should be something that produces a derivation, not
# a string or a raw attribute set or anything else.
# buildGoModule is a function that returns a derivation, so
# you want `buildGoModule ...` here, not `{ pet = ...; }` here;
# the latter is an attribute set.

buildGoModule rec {
  pname = "windowarranger";
  version = "0.1";

  src = ./.;
#  src = fetchFromGitHub {
#    owner = "surlykke";
#    repo = "";
#    rev = "v${version}";
#    hash = "sha256-Gjw1dRrgM8D3G7v6WIM2+50r4HmTXvx0Xxme2fH9TlQ=";
#  };

	#  buildInputs = [ pkg-config ] ++ libs;

  # this hash is updated from the example, which seems to be out of date
	#  vendorHash = lib.fakeHash;
	vendorHash = null;

  meta = {
    description = "Window arranger for sway";
    homepage = "https://github.com/surlykke/windowarranger";
    license = lib.licenses.gpl2;
    maintainers = with lib.maintainers; [ surlykke ];
  };
}
