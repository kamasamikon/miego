
FROM alpine
# FROM centos:7
# FROM frolvlad/alpine-python3
# FROM openresty/openresty

WORKDIR /root

apk update
apk add nginx

COPY nginx.conf.templ /etc/nginx/
COPY dockerstuff/openresty/nginx/conf/server.crt /etc/nginx
COPY dockerstuff/openresty/nginx/conf/server.csr /etc/nginx
COPY dockerstuff/openresty/nginx/conf/server.key  /etc/nginx
COPY dockerstuff/openresty/nginx/conf/server_nopwd.key  /etc/nginx


### #####################################################################
## Install MIE
#
COPY ./PRC /etc/localtime


### #####################################################################
## Run
#
COPY dockerentry. .
ENTRYPOINT /root/dockerentry.sh
