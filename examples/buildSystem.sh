#!/bin/bash

firmwareName=anytool

echo "building"
go build -ldflags "-s -w" -o ${firmwareName} .
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
