#!/bin/sh -e

/usr/local/bin/mount_home_dirs.sh &

/usr/local/bin/docker-init /usr/local/bin/dockerd -- \
  --config-file /var/config/docker/daemon.json \
  --swarm-default-advertise-addr=eth0 \
  --userland-proxy-path /usr/bin/vpnkit-expose-port \
  --storage-driver overlay2
