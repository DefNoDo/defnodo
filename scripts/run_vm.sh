#!/usr/bin/env bash
set -e

#./output/linuxkit run hyperkit -fw bunk_uefi.fd -hyperkit output/hyperkit -vpnkit output/vpnkit -cpus 4 -mem 8192 -disk size=10G -networking=vpnkit  -vsock-ports 2376 -squashfs -data-file metadata.json vm-image/defnodo
./output/linuxkit run hyperkit -fw bunk_uefi.fd -hyperkit output/hyperkit -vpnkit output/vpnkit -cpus 4 -mem 8192 -disk size=10G -networking=vpnkit  -vsock-ports 2376 -squashfs -data-file metadata.json defnodo-data/defnodo

# ./output/hyperkit -A -u -F vm-image/defnodo-state/hyperkit.pid -c 1 -m 4096M \
#   -s 0:0,hostbridge -s 31,lpc \
#   -s 1:0,virtio-vpnkit,path=vm-image/defnodo-state/vpnkit_eth.sock,uuid=5b4b1f75-7176-42a7-b9ba-bafbafd98a61 -U 897095f6-adb4-471d-9a0c-3c5273b9484a \
#   -s 2:0,virtio-blk,vm-image/defnodo-squashfs.img \
#   -s 3:0,ahci-hd,vm-image/defnodo-state/disk.raw \
#   -s 4,virtio-sock,guest_cid=3,path=vm-image/defnodo-state,guest_forwards="2376;62373" \
#   -s 5,ahci-cd,vm-image/defnodo-state/data.iso \
#   -s 6,virtio-rnd \
#   -s 7,virtio-9p,path=vm-image/defnodo-state/vpnkit_port.sock,tag=port \
#   -s 8,virtio-9p,path=vm-image/defnodo-state/mount_home_9p.sock,tag=home_mount \
#   -l com1,stdio,log=vm-image/defnodo-state/console-ring -f kexec,vm-image/defnodo-kernel,,earlyprintk=serial console=ttyS0 page_poison=1 root=/dev/vda

# ./output/hyperkit -A -u -F vm-image/defnodo-state/hyperkit.pid -c 1 -m 1024M \
#   -s 0:0,hostbridge \
#   -s 31,lpc \
#   -s 1:0,virtio-vpnkit,path=vm-image/defnodo-state/vpnkit_eth.sock,uuid=57eb5a59-dac1-4410-a792-e391c13c7136 -U 6bfd42c3-7d2e-4812-96cb-f4ff7a184c7b \
#   -s 2:0,virtio-blk,vm-image/defnodo-squashfs.img \
#   -s 3:0,ahci-hd,vm-image/defnodo-state/disk.raw \
#   -s 4,virtio-sock,guest_cid=3,path=vm-image/defnodo-state,guest_forwards="2376;62373" \
#   -s 5,ahci-cd,vm-image/defnodo-state/data.iso \
#   -s 6,virtio-rnd \
#   -s 7,virtio-9p,path=vm-image/defnodo-state/vpnkit_port.sock,tag=port \
#   -l com1,stdio,log=vm-image/defnodo-state/console-ring -f kexec,vm-image/defnodo-kernel,,earlyprintk=serial console=ttyS0 page_poison=1 root=/dev/vda
# output/vpnkit --ethernet vm-image/defnodo-state/vpnkit_eth.sock --vsock-path vm-image/defnodo-state/connect --port vm-image/defnodo-state/vpnkit_port.sock
