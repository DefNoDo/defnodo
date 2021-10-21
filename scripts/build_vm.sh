#!/usr/bin/env bash
set -e

mkdir -p vm-image/

./output/linuxkit -v build -format "kernel+squashfs" -dir vm-image/ linuxkit/defnodo.yml
