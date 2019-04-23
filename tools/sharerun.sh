#!/bin/sh

set -v

NAME=$1
LOCALDIR=$2
FORCE=$3

echo "RUN $LOCALDIR/main BY MSA."

MSBIP=$(sudo docker inspect --format '{{ .NetworkSettings.IPAddress }}' msb)

if [ "$FORCE" != "" ]; then
    sudo docker rm -f "${NAME}"
fi

sudo docker run -it --name ${NAME} \
    -v /tmp/.conf.${NAME}:/tmp/conf \
    -v ${LOCALDIR}:/root/ms \
    -e MSBHOST=${MSBIP} \
    msa
