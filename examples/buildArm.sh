#!/bin/bash

firmwareName=anytool-linux-arm5

echo "building"
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=5 CC=arm-linux-gcc go build -ldflags "-s -w" -o ${firmwareName} .
#xgo --targets=linux/arm-5 -ldflags "-s -w" .
if [ $? -ne 0 ]
then
	echo "build failed"
	exit
fi

bzip2 -c ${firmwareName} > ${firmwareName}.bz2
if [ $? -eq 0 ]
then
    echo "build success"
else
    echo "build failed"
fi
