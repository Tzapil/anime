#!/bin/bash

rm -r ./bin
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.8 ./build_docker.sh

docker build -t tzapil/anime:v0.3 -f Dockerfile .