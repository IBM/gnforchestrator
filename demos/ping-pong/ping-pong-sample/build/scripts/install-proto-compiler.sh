#!/usr/bin/env bash

if which protoc > /dev/null; then
    exit 0;
fi

# Make sure you grab the latest version
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.11.2/protoc-3.11.2-linux-x86_64.zip

# Unzip
unzip protoc-3.11.2-linux-x86_64.zip -d protoc3

# Move protoc to /usr/local/bin/
sudo mv protoc3/bin/* /usr/local/bin/

# Move protoc3/include to /usr/local/include/
sudo mv protoc3/include/* /usr/local/include/

rm protoc-3.11.2-linux-x86_64.zip
rm -r  protoc3/

go get -u github.com/golang/protobuf/protoc-gen-go