#!/bin/sh

./makesrc.sh
sudo docker rm -f msb
sudo docker rmi -f msb
sudo docker build -t msb .
