#!/bin/sh

NAME=$1
LOCALDIR=$2

echo "RUN $LOCALDIR/main BY MSA."

echo "msabyshare NAME LOCALDIR"
echo "sudo docker run -it --name \$NAME -v \$LOCALDIR:/service msa"

sudo docker rm -f ${NAME}
sudo docker run -it --name ${NAME} -v ${LOCALDIR}:/service msa
