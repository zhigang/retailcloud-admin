#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

echo "go build..."
# CGO_ENABLED=0 GOARCH=amd64 go build -o build/app
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/app
# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/app.exe

echo "copy config..."
cp config/config.yml build/

docker build -t registry.cn-zhangjiakou.aliyuncs.com/ys-oms-test/retailcloud-tool:1.2 .

echo "build successful!"