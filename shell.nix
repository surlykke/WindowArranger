let
  pkgs = import <nixpkgs> {};
in
pkgs.mkShell {
	buildInputs = with pkgs; [
		go
		gopls
		gotools
		go-tools
	];
}

