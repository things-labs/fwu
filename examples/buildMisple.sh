#!/bin/bash

firmwareName=anytool-linux-mipsle

echo "building"
CGO_ENABLED=1 GOOS=linux GOARCH=mipsle \
STAGING_DIR=/opt/toolchain/openwrt18.06/staging_dir \
CC=/opt/toolchain/openwrt18.06/staging_dir/gcc-mipsel-linux-7.3.0/bin/mipsel-openwrt-linux-gcc \
go build -ldflags "-s -w" -o ${firmwareName} .

if [ $? -ne 0 ]
then
	echo "xgo failed"
	exit
fi

bzip2 -c ${firmwareName} > ${firmwareName}.bz2
if [ $? -eq 0 ]
then
    echo "build success"
else
    echo "build failed"
fi
