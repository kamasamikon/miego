#!/bin/sh 

sudo docker rm --force msb
sudo docker run -it -d --name msb -p 9080:80 msb:latest
