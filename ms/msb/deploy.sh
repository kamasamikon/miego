#!/bin/bash

REMOTE=root@117.50.95.113
REMOTE=root@106.75.49.211

sudo docker save -o msb.tar msb

sudo rm msb.tar.gz
sudo chmod a+rwx msb.tar
sudo gzip --best -c msb.tar >> msb.tar.gz
sudo scp msb.tar.gz $REMOTE:/tmp/msb.tar.gz

echo ">>>>>>>>>>>>>>" ssh $REMOTE "docker rm -f msb"
ssh $REMOTE "docker rm -f msb"

echo ">>>>>>>>>>>>>>" ssh $REMOTE "docker load -i /tmp/msb.tar"
ssh $REMOTE "gunzip /tmp/msb.tar.gz"
ssh $REMOTE "docker load -i /tmp/msb.tar"

echo ">>>>>>>>>>>>>>" ssh $REMOTE "docker run -d -it --restart=always --name msb -p 9080:80 msb"

ssh $REMOTE "docker rm -f msb"
ssh $REMOTE "docker run -d -it --restart=always --name msb -p 9080:80 msb"

ssh $REMOTE "docker rm -f msb.mp"
ssh $REMOTE "docker run -d -it --restart=always --name msb.mp -p 9081:80 msb"

