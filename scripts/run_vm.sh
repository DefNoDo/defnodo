#!/usr/bin/env bash
set -e

# ./linuxkit run hyperkit -fw bunk_uefi.fd  -disk size=4G -networking=vpnkit  -vsock-ports 2376 -squashfs -data-file metadata.json -mem 4096 docker-for-mac
./output/linuxkit run hyperkit -fw bunk_uefi.fd  -disk size=4G -networking=vpnkit  -vsock-ports 2376 -squashfs -data-file metadata.json -mem 4096 vm-image/defnodo