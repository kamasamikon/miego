#!/bin/sh

set -v

NAME=$1
LOCALDIR=$2

echo "RUN $LOCALDIR/main BY MSA."

MSBIP = `sudo docker inspect --format '{{ .NetworkSettings.IPAddress }}' msb`

sudo docker rm -f ${NAME}

sudo docker run -it --name ${NAME} \
    -v /tmp/.conf.${NAME}:/tmp/conf \
    -v ${LOCALDIR}:/root/ms \
    -e MSBHOST=${MSBIP} \
    msa
