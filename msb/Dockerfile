
FROM nginx:1.19.6-alpine
# FROM smebberson/alpine-nginx
# FROM alpine
# FROM openresty/openresty

WORKDIR /root

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories


# RUN apk update \
# && apk upgrade \
# && apk add --no-cache nginx 

COPY nginx.conf.templ /etc/nginx/
COPY nginx.conf.templ /etc/nginx/nginx.conf

#COPY dockerstuff/openresty/nginx/conf/server.crt /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server.csr /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server.key  /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server_nopwd.key  /etc/nginx


### #####################################################################
## Install MIE
#
COPY ./PRC /etc/localtime
COPY ./msb /root/msb

VOLUME /var/lib/nginx/tmp/client_body

### #####################################################################
## Run
#
COPY dockerentry.sh .
ENTRYPOINT /root/dockerentry.sh
