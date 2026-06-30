FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w -buildid=" -o /out/md2img .

FROM alpine:3
RUN apk upgrade --no-cache --available && \
    apk add --no-cache \
      chromium \
      fontconfig \
      font-noto-cjk \
      font-noto-emoji \
    && rm -rf /var/cache/apk/*
COPY manifest/deploy/local.conf /etc/fonts/local.conf
RUN fc-cache -f
RUN adduser -D chrome
COPY --from=builder /out/md2img /app/md2img
ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/ \
    CHROMIUM_FLAGS="--disable-software-rasterizer --disable-dev-shm-usage"
USER chrome
WORKDIR /app
ENTRYPOINT ["/app/md2img"]
