FROM alpine
# FROM smebberson/alpine-nginx

WORKDIR /msa

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.4/main/" > /etc/apk/repositories

# RUN apk update && apk upgrade && apk add --no-cache nginx 

COPY ./PRC /etc/localtime
COPY ./msa.cfg /msa/msa.cfg
COPY ./msa /msa

RUN mkdir /service
ENTRYPOINT /msa/msa
