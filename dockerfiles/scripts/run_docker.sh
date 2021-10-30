#!/bin/sh -e

/usr/local/bin/mount_home_dirs.sh &
/usr/local/bin/containerd &
/usr/local/bin/docker-init /usr/local/bin/dockerd -- \
  --config-file /var/config/docker/daemon.json \
  --swarm-default-advertise-addr=eth0 \
  --containerd /run/containerd2/containerd.sock \
  --userland-proxy-path /usr/bin/vpnkit-expose-port \
  --storage-driver overlay2

# Previous/in yaml file definition
# "/usr/local/bin/docker-init", "/usr/local/bin/dockerd", "--",
#             "--config-file", "/var/config/docker/daemon.json",
#             "--swarm-default-advertise-addr=eth0",
#             # "--containerd", "/var/run/containerd/containerd.sock",
#             "--userland-proxy-path", "/usr/bin/vpnkit-expose-port",
#             "--storage-driver", "overlay2"