FROM zenika/alpine-chrome:latest

USER root

# 更新并安装中文字体
RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories && \
    apk update && apk add --no-cache \
        ttf-dejavu \
        fontconfig \
        && \
    fc-cache -f  && \
    rm -rf /var/cache/apk/*

USER chrome

COPY md2img /app/md2img

WORKDIR /app

# 使用 ENTRYPOINT 来支持传参
ENTRYPOINT ["/app/md2img"]
