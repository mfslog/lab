#!/usr/bin/env bash
protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf \
          --gofast_out=../idl/ *.proto