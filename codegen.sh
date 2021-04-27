#!/usr/bin/env bash

rm -rf ./ports/grpc/proto/*.pb.go
protoc --proto_path=ports/grpc/proto ./ports/grpc/proto/*.proto --go_out=plugins=grpc:ports/grpc/proto