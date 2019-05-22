FROM alpine

WORKDIR /root

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/latest-stable/main/" > /etc/apk/repositories \
        && echo "https://mirror.tuna.tsinghua.edu.cn/alpine/latest-stable/community/" >> /etc/apk/repositories \
        && apk update \
        && apk upgrade \
        && apk add redis \
        && rm -rf /var/cache/apk/* 

COPY ./PRC /etc/localtime
COPY ./msa .

RUN mkdir /ms
ENTRYPOINT /root/msa
