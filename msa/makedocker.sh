#!/bin/sh

./makesrc.sh
upx msa
sudo docker build -t msa .
