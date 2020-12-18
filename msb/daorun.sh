#!/bin/sh 

./makedocker.sh

sudo docker rm --force msb

echo ">>>> " sudo docker run --restart=always -it $@ --name msb -p 9080:80 msb:latest
sudo docker run --restart=always -it $@ --name msb -p 9080:80 msb.slim:latest
