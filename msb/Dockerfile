
FROM smebberson/alpine-nginx
# FROM alpine
# FROM centos:7
# FROM frolvlad/alpine-python3
# FROM openresty/openresty

WORKDIR /root

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories


# RUN apk update \
# && apk upgrade \
# && apk add --no-cache nginx 

COPY nginx.conf.templ /etc/nginx/
#COPY dockerstuff/openresty/nginx/conf/server.crt /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server.csr /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server.key  /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server_nopwd.key  /etc/nginx


### #####################################################################
## Install MIE
#
COPY ./PRC /etc/localtime
COPY ./main/main /root/main


### #####################################################################
## Run
#
COPY dockerentry.sh .
ENTRYPOINT /root/dockerentry.sh
