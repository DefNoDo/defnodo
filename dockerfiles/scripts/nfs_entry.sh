#! /bin/sh -e

echo "MOUNTPOINT: ${MOUNTPOINT}"
echo "SERVER: ${SERVER}"
echo "SHARE: ${SHARE}"
echo "MOUNT_OPTIONS: ${MOUNT_OPTIONS}"

mkdir -p "$MOUNTPOINT"
ls -l /host_var
ls -l ${MOUNTPOINT}

rpc.statd & rpcbind -f &
mount -t nfs -o port=7777,mountport=7777,nfsvers=3,noacl,tcp 192.168.65.2:/mount /host_var/nfs
# mount -t "$FSTYPE" -o "$MOUNT_OPTIONS" "$SERVER:$SHARE" "$MOUNTPOINT"
mount | grep nfs

echo "Mounted $MOUNTPOINT"
ls -l $MOUNTPOINT