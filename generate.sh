#!/bin/bash

protoc pubsubpb/pubsub.proto --go_out=plugins=grpc:.