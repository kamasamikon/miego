FROM alpine

WORKDIR /root

# RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories

# RUN apk update && apk upgrade && apk add --no-cache nginx 

COPY ./PRC /etc/localtime
COPY ./msa .

RUN mkdir /ms
ENTRYPOINT /root/msa
