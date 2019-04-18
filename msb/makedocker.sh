#!/bin/sh

./makesrc.sh
sudo docker build -t msb .
