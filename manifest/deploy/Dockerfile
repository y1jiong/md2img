FROM selenium/standalone-chrome:latest

COPY md2img /app/md2img

WORKDIR /app

# 使用 ENTRYPOINT 来支持传参
ENTRYPOINT ["/app/md2img"]
