#!/bin/sh 

docker rm --force msb
docker run -it -d --name msb -p 9080:80 msb:latest
