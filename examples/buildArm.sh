#!/bin/bash

xgo --targets=linux/arm-5 -ldflags "-s -w" .
#go build -ldflags "-s -w"
if [ $? -ne 0 ]
then
	echo "xgo failed"
	exit
fi

bzip2 -c examples-linux-arm-5 > anytool.bz2
#bzip2 -c examples > anytool.bz2
if [ $? -eq 0 ]
then
    echo "build success"
fi
