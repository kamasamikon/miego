#!/bin/sh

./makesrc.sh
sudo docker build -t msa .
