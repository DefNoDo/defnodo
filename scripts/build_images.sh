#!/usr/bin/env bash
set -e

function finish {
  popd
}

trap finish EXIT

pushd dockerfiles

echo "Building NFS client..."
docker build -t defnodo/nfs-client:latest . -f Dockerfile.nfsclient

echo "Building vpnkit-forwarder..."
docker build -t defnodo/vpnkit-forwarder:v0.5.0-custom . -f Dockerfile.vpnkit-forwarder