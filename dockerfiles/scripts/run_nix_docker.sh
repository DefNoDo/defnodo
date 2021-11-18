#!/bin/bash -e

export HOME=/home/nixuser
export USER=nixuser
. /home/nixuser/.nix-profile/etc/profile.d/nix.sh
unset USER

/usr/local/bin/mount_home_dirs.sh &

MOUNT_DIR=$(cat /var/config/docker/user.config | tr -d '\n')
DOCKER_VERSION=$(cat /var/config/docker/docker.version | tr -d '\n' | awk '{print $1}')
DOCKER_HASH=$(cat /var/config/docker/docker.version | tr -d '\n' | awk '{print $2}')

if [ $DOCKER_VERSION == "latest" ]; then
  dockerd --version || nix-env -i docker
else
  dockerd --version | grep -q "${DOCKER_VERSION}" || nix-env -i docker-${DOCKER_VERSION} -f https://github.com/NixOS/nixpkgs/archive/${DOCKER_HASH}.tar.gz
fi

dockerd \
  --config-file /var/config/docker/daemon.json \
  --swarm-default-advertise-addr=eth0 \
  --userland-proxy-path /usr/bin/vpnkit-expose-port \
  --storage-driver overlay2

# Podman notes:
# nix-env -i podman runc conmon skopeo slirp4netns fuse-overlayfs