{ pkgs ? import <nixos>{}}:
let
  drv = { stdenv, fetchFromGitHub, cmake }:
stdenv.mkDerivation rec {
  pname = "NGT";
  version = "v1.9.0";
  nativeBuildInputs = [ cmake ];
  buildInputs = [ ];
  NIX_ENFORCE_NO_NATIVE=false;
  __AVX2__=1;
  src = /home/tom/NGT;
  # src = fetchFromGitHub {
  #   owner = "yahoojapan";
  #   repo = "NGT";
  #   rev = version;
  #   sha256 = "19796e45c6921be8e1ebeadc75f4d648600e3a505337484ff17966350b7913cc";
  # };
  enableParallelBuilding = true;
};


  ngt = pkgs.callPackage drv {};
in
pkgs.mkShell {
  buildInputs = with pkgs; [
    hdf5
    ngt
  ];
  nativeBuildInputs = with pkgs; [ pkgconfig gcc gnumake go_1_13 ] ;
  dontDisableStatic = true;
  GOPATH="/home/tom/go";
  shellHook = ''
    '';
  
}
