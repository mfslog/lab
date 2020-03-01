#!/usr/bin/env bash

protoc -I=. \
       --go_out=../proto ./account/acct_def.proto

protoc -I=. -I=./account\
       --go_out=plugins=grpc,Macct_def.proto=github.com/mfslog/lab/go/kit/proto/account:../proto ./account/acct.proto

