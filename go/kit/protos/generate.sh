#!/usr/bin/env bash

protoc -I=. \
       --go_out=../idl ./account/acct_def.proto

protoc -I=. -I=./account\
       --go_out=plugins=grpc,Macct_def.proto=github.com/mfslog/lab/kit/idl/account:../idl ./account/acct.proto

