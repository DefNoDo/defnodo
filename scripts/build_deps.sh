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

function build_hyperkit {
  make all
  cp build/hyperkit ../../output/hyperkit
}

function build_go9p {
  go build -o ../../output/go9p cmd/export9p/main.go
}

function get_vpnkit_tools {
  HOMEBREW_NO_AUTO_UPDATE=1
  brew install wget pkg-config dylibbundler libtool automake
  brew install opam
  opam env || opam init --compiler 4.12.0 -n
}

if [ $# -eq 0 ]; then
  set - "linuxkit" "vpnkit" "hyperkit" "go9p"
fi

while (($#)); do
  name=${1}
  command="build_${name}"
  echo "Building ${name}..."
  if [ $1 == "vpnkit" ]; then
    get_vpnkit_tools
  fi

  (cd deps/${name}/ && ${command})
  shift
done
