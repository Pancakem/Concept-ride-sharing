#!/bin/bash


cd rides/v1/src/cmd/app
go get -d ./...

docker pull redis
docker run --name some-redis -d -p 6379 redis

mkdir keys

cd keys

ssh-keygen -t rsa -N "" -f key
