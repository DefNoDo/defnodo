#!/usr/bin/env bash
set -e

DOCKER_HOST=""
if [ $# -gt 0 ]; then
  DOCKER_HOST="${1}"
fi

export DOCKER_HOST

echo "Building vpnkit-forwarder..."
docker build -t defnodo/vpnkit-forwarder:v0.5.0-custom dockerfiles/. -f dockerfiles/Dockerfile.vpnkit-forwarder

echo "Building defnodo utils..."
docker build -t defnodo/utils:latest . -f dockerfiles/Dockerfile.utils