#!/usr/bin/env bash
set -e

function finish {
  popd
}

trap finish EXIT

pushd deps

echo "Building linuxkit..."
pushd linuxkit
make all
cp bin/linuxkit ../../output/linuxkit
popd