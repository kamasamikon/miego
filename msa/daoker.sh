#!/bin/sh 


#
# THIS IS A DEMO
# THIS IS A DEMO
# THIS IS A DEMO
# THIS IS A DEMO
# THIS IS A DEMO
# THIS IS A DEMO
# THIS IS A DEMO
#

# MS information
msName=vcode
msVern=v1
msPort=8888
msDesc="Generate QR code for given string content."

# User's Docker file
cat << __EOF__ > dfUser
RUN echo what's up.
__EOF__

# userScript called before build
cat << __EOF__ > userScript

# Make/Build
CGO_ENABLED=0 go build -ldflags "-w -s" main.go

cp -frv templates ms

__EOF__

# Make the docker
daoker.py -n $msName -v $msVern -p $msPort -d "$msDesc" -D ./dfUser $@ 

# Cleanup
echo
echo
echo "------------------------------------------------------------"
cat Dockerfile
echo "------------------------------------------------------------"

rm Dockerfile
rm dfUser
rm userScript
