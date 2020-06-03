#!/bin/bash

REMOTE=root@106.75.49.211

sudo docker save -o msb.tar msb

sudo chmod a+rwx msb.tar
sudo scp msb.tar $REMOTE:/tmp/msb.tar 

echo ">>>>>>>>>>>>>>" ssh $REMOTE "docker rm -f msb"
ssh $REMOTE "docker rm -f msb"

echo ">>>>>>>>>>>>>>" ssh $REMOTE "docker load -i /tmp/msb.tar"
ssh $REMOTE "docker load -i /tmp/msb.tar"

echo ">>>>>>>>>>>>>>" ssh $REMOTE "docker run -d -it --restart=always --name msb -p 9080:80 msb"
ssh $REMOTE "docker run -d -it --restart=always --name msb -p 9080:80 msb"

