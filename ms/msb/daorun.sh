#!/bin/sh -x

./makedocker.sh

cName=msb

#iName=msb.slim
iName=msb

sudo docker rm --force ${cName}
sudo docker rm --force ${cName}.mp

echo ">>>> " sudo docker run --restart=always -it "$@" --name $cName -p 9080:80 $iName
sudo docker run --restart=always -d -it "$@" --name ${cName} -p 9080:80 -p 9443:443 $iName
# sudo docker run --restart=always -d -it "$@" --name ${cName}.mp -p 9081:80 -p 10443:443 $iName
