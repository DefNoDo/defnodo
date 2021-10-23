#!/usr/bin/env bash
set -ex

function build_linuxkit {
  make all
  cp bin/linuxkit ../../output/linuxkit
}

function build_vpnkit {
  OPAMYES="1"
  OPAMJOBS="2"
  OPAMVERBOSE=1
  eval $(opam env --switch=4.12.0)
  echo "Building vpnkit ocaml requirements..."
  make -f Makefile.darwin ocaml
  echo "Building vpnkit dependencies..."
  make -f Makefile.darwin depends
  echo "Building vpnkit executable..."
  make -f Makefile.darwin build
  cp _build/default/src/bin/main.exe ../../output/vpnkit
}

function get_vpnkit_tools {
  HOMEBREW_NO_AUTO_UPDATE=1
  brew install wget pkg-config dylibbundler libtool automake
  brew install opam
  opam env || opam init --compiler 4.12.0 -n

}

echo "Building linuxkit..."
(cd deps/linuxkit/ && build_linuxkit)

echo "Building vpnkit..."
get_vpnkit_tools
(cd deps/vpnkit/ && build_vpnkit $(pwd))
