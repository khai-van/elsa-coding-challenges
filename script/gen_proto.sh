#!/bin/bash
protoc --go_out=api/gen --go-grpc_out=api/gen --proto_path=api/proto api/proto/*.proto