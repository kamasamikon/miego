#!/bin/sh

rm /root/msb.cfg

function setcfg() {
    echo "$1=$2" >> /root/msb.cfg
}

if [ -x /usr/local/openresty/bin/openresty ]; then
    # /root/msb.cfg
    setcfg s:/msb/nginx/conf /etc/nginx/conf.d/msb.conf
    setcfg s:/msb/nginx/tmpl /etc/nginx/conf.d/msb.conf.tmpl
    setcfg s:/msb/nginx/exec /usr/local/openresty/bin/openresty

    # copy nginx files, use grpc.tmpl
    rm /etc/nginx/conf.d/default.conf
    cp /root/nginx.conf.grpc.tmpl /etc/nginx/conf.d/msb.conf
    cp /root/nginx.conf.grpc.tmpl /etc/nginx/conf.d/msb.conf.tmpl

    # run
    /usr/local/openresty/bin/openresty &
else
    # /root/msb.cfg
    setcfg s:/msb/nginx/conf /etc/nginx/nginx.conf
    setcfg s:/msb/nginx/tmpl /etc/nginx/nginx.conf.tmpl
    setcfg s:/msb/nginx/exec /usr/sbin/nginx

    # copy nginx files, use http.tmpl
    cp /root/nginx.conf.http.tmpl /etc/nginx/conf.d/msb.conf
    cp /root/nginx.conf.http.tmpl /etc/nginx/conf.d/msb.conf.tmpl

    # run
    nginx &
fi


while true; do 
    cd /root
    cat /root/msb.cfg
    /root/msb 
done
