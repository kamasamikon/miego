
# XXX: For grpc, use openresty/openresty:alpine, else use nginx:alpine
# FROM openresty/openresty:alpine
FROM nginx:1.19.6-alpine
# FROM smebberson/alpine-nginx

WORKDIR /root

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories

COPY nginx.conf.http.tmpl .
COPY nginx.conf.grpc.tmpl .

COPY msb.cfg /etc/msb.cfg

COPY msb.key /etc/nginx
COPY msb.pem /etc/nginx

#COPY dockerstuff/openresty/nginx/conf/server.crt /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server.csr /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server.key  /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server_nopwd.key  /etc/nginx

COPY ./PRC /etc/localtime
COPY ./msb .

COPY ./now .

VOLUME /var/lib/nginx/tmp/client_body

ENV KLOG_MUTE=0

### #####################################################################
## Run
#
COPY dockerentry.sh .
ENTRYPOINT ["/root/dockerentry.sh"]
