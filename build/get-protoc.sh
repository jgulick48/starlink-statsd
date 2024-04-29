#!/bin/bash

# requires https://github.com/fullstorydev/grpcurl
go run github.com/fullstorydev/grpcurl/cmd/grpcurl@latest -plaintext -protoset-out proto/dish.protoset 192.168.100.1:9200 describe SpaceX.API.Device.Device