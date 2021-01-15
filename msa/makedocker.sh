#!/bin/sh

Type=$1

./makesrc.sh
upx msa
echo sudo docker build -f Dockerfile.$Type -t msa .
sudo docker build -f Dockerfile.$Type -t msa-$Type .
