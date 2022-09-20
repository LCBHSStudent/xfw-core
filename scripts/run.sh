#!/bin/bash

pushd /xfw/xfw-core/bin

service mysql start
nohup ./go-cqhttp > cqhttp.log &2>1 &
sleep 10 
nohup ./xfw-core > xfw.log &2>1 &
nohup ./homo-space > homo.log &2>1 &

popd

while [[ true ]]; do
    sleep 1
done
