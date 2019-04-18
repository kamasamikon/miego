#!/bin/sh 

./makedocker.sh

sudo docker rm --force msb
sudo docker run -it $@ --name msb -p 9080:80 msb:latest
