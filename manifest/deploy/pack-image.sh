#!/bin/bash

USER=y1jiong
BIN=md2img
GIT_TAG=$(shell git describe --tags --abbrev=0)

echo "> tar -xJvf ${BIN}.*.tar.xz"
tar -xJvf ${BIN}.*.tar.xz
if [ $? -ne 0 ]; then
    echo "failed"
    exit 1
fi

echo "> chmod 700 ${BIN}"
chmod 700 ${BIN}

echo "docker build -t ${USER}/${BIN}:${GIT_TAG} ."
docker build -t ${USER}/${BIN}:${GIT_TAG} .

echo "> rm ${BIN}.*.tar.xz"
rm ${BIN}.*.tar.xz
