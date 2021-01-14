#!/bin/sh 

./makedocker.sh

cName=msb.grpc

#iName=msb.slim
iName=msb

sudo docker rm --force $cName

echo ">>>> " sudo docker run --restart=always -it $@ --name $cName -p 9080:80 $iName
sudo docker run --restart=always -it $@ --name $cName -p 9080:80 $iName
