#!/bin/bash

#go build -o ./bin/main .
rm -r ./bin
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main .

docker build -t tzapil/anime:v0.2 -f Dockerfile.scratch .