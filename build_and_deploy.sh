#!/usr/bin/bash

echo "Checking GOPATH..."

echo ${GOPATH:?"-> not set, exitting..."}

cd $GOPATH

export GOOS=linux
export GOARCH=386

go build -v github.com/brockwood/govantage > build.log 2>&1

if [ "$?" -ne "0" ]; then
	echo "Error building the executable, please check the build log..."
else
	echo "Build was successful, deploying to the Edison..."
fi

scp $GOPATH/govantage root@rock-iot:~

if [ "$?" -ne "0" ]; then
	echo "Error copying the executable to the remote Edison..."
else
	echo "Build was copied to the remote Edison."
fi

ssh root@rock-iot "chmod +x /home/root/govantage"