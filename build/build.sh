#!/bin/bash
set -e
self_dir=$(cd $(dirname $0); pwd)
app_dir=$(dirname $self_dir)

cd $app_dir

IMAGE_NAME=utsso/certmon
VERSION=$(date +%y.%m.%d)

GOOS=linux CGO_ENABLED=0 go build -o debug/certmon ./cmd

docker build -t $IMAGE_NAME:$VERSION -f ./build/Dockerfile .
docker tag $IMAGE_NAME:$VERSION $IMAGE_NAME:latest

if [[ X"$1" == X"prod" ]]; then
    docker push $IMAGE_NAME:$VERSION
    docker push $IMAGE_NAME:latest
fi