#!/usr/bin/env bash

if ! which protoc > /dev/null; then
    echo "error: protoc not installed, please run make proto-compiler" >&2;
    exit 1;
fi

protoc --go_out=plugins=grpc:. ./protos/*.proto