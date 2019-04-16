
FROM alpine
# FROM smebberson/alpine-nginx
# FROM centos:7
# FROM frolvlad/alpine-python3
# FROM openresty/openresty

WORKDIR /root

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories


# RUN apk update && apk upgrade && apk add --no-cache nginx 


### #####################################################################
## Install RPC
#
COPY ./PRC /etc/localtime

### #####################################################################
## Run
#
COPY ./msa .
COPY dockerentry.sh .
ENTRYPOINT /root/dockerentry.sh
