FROM zenika/alpine-chrome:latest

USER root

# 更新并安装中文字体（文泉驿不好看还得是思源）
RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories && \
    apk update && apk add --no-cache \
        fontconfig \
        font-noto-cjk \
        && \
    apk del font-wqy-zenhei && \
    rm -rf /var/cache/apk/*

COPY local.conf /etc/fonts/local.conf
RUN fc-cache -f

USER chrome

COPY md2img /app/md2img
WORKDIR /app

# 使用 ENTRYPOINT 来支持传参
ENTRYPOINT ["/app/md2img"]
