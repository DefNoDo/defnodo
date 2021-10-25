#!/bin/sh -e

echo "Mounting port..."
mkdir -p /host_var/vpnkit
mount -v -t 9p -o trans=virtio,dfltuid=1001,dfltgid=50,version=9p2000 port /host_var/vpnkit

echo "Mounting homedir..."
mkdir -p /host_var/home_mount
mount -v -t 9p -o trans=tcp,port=7777,access=any,dfltuid=501,dfltgid=50,version=9p2000 192.168.65.2 /host_var/home_mount
