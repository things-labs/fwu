#!/bin/bash

echo "building"
go build -ldflags "-s -w" -o anytool .
if [ $? -ne 0 ]
then
	echo "build failed"
	exit
fi

bzip2 -c anytool > anytool.bz2
if [ $? -eq 0 ]
then
    echo "build success"
else
    echo "build failed"
fi
