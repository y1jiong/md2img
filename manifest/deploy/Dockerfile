FROM alpine:latest

# Configure Alpine mirror (optional, for faster downloads in China)
RUN sed -i 's#https\?://dl-cdn.alpinelinux.org/alpine#https://mirrors.tuna.tsinghua.edu.cn/alpine#g' /etc/apk/repositories

# Install Chromium and fonts in a single layer
RUN apk upgrade --no-cache --available && \
    apk add --no-cache \
      chromium-swiftshader \
      fontconfig \
      font-noto-cjk \
      font-noto-emoji \
    && \
    rm -rf /var/cache/apk/*

# Add font configuration
COPY local.conf /etc/fonts/local.conf
RUN fc-cache -f

# Add Chrome as a user
RUN mkdir -p /usr/src/app && \
    adduser -D chrome && \
    chown -R chrome:chrome /usr/src/app

# Copy application
COPY --chown=chrome:chrome md2img /app/md2img

# Set environment variables
ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/ \
    CHROMIUM_FLAGS="--disable-software-rasterizer --disable-dev-shm-usage"

# Switch to non-privileged user
USER chrome
WORKDIR /app

# Set the entrypoint
ENTRYPOINT ["/app/md2img"]