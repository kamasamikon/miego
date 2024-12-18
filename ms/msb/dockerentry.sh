#!/bin/sh

[ -f /root/msb.cfg ] && rm /root/msb.cfg
rm /etc/nginx/conf.d/default.conf

setcfg() {
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
    setcfg s:/msb/nginx/conf /etc/nginx/conf.d/msb.conf
    setcfg s:/msb/nginx/tmpl /etc/nginx/conf.d/msb.conf.tmpl
    setcfg s:/msb/nginx/exec /usr/sbin/nginx

    # copy nginx files, use http.tmpl
    cp /root/nginx.conf.http.tmpl /etc/nginx/conf.d/msb.conf
    cp /root/nginx.conf.http.tmpl /etc/nginx/conf.d/msb.conf.tmpl

    # run
    nginx &
fi


while true; do 
    cd /root || return
    echo ">>> cat /root/msb.cfg"
    cat /root/msb.cfg
    echo "<<< cat /root/msb.cfg"
    /root/msb 
done
