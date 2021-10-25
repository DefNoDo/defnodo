#!/bin/sh -e

MOUNT_DIR="/Users/mbrace"
apk add fuse

mkdir -p "${MOUNT_DIR}"
/usr/local/bin/mount9p_fuse -v 192.168.65.2:7777 "${MOUNT_DIR}" &

while 1; do
  ls "${MOUNT_DIR}" > /dev/null
done