
FROM nginx:1.19.6-alpine
# FROM smebberson/alpine-nginx

WORKDIR /root

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories

# No neccessary to copy nginx.conf.server.templ
COPY nginx.conf.full.templ /etc/nginx/nginx.conf.templ
COPY nginx.conf.full.templ /etc/nginx/nginx.conf

COPY msb.key /etc/nginx
COPY msb.pem /etc/nginx

#COPY dockerstuff/openresty/nginx/conf/server.crt /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server.csr /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server.key  /etc/nginx
#COPY dockerstuff/openresty/nginx/conf/server_nopwd.key  /etc/nginx

COPY ./PRC /etc/localtime
COPY ./msb /root/msb

VOLUME /var/lib/nginx/tmp/client_body

### #####################################################################
## Run
#
COPY dockerentry.sh .
ENTRYPOINT /root/dockerentry.sh
