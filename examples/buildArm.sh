#!/bin/bash

CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=5 CC=arm-linux-gcc go build -ldflags "-s -w" -o anytool-linux-arm-5 .
#xgo --targets=linux/arm-5 -ldflags "-s -w" .
if [ $? -ne 0 ]
then
	echo "xgo failed"
	exit
fi

bzip2 -c anytool-linux-arm-5 > anytool-arm5.bz2
if [ $? -eq 0 ]
then
    echo "build success"
fi
