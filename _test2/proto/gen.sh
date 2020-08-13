#!/usr/bin/env bash

# generate the gRPC code
protoc -I/usr/local/include -I. ${GOPATHLIST} --go_out=plugins=grpc:. \
    hello.proto
