#!/bin/bash

export GOPATH=$GOPATH:`pwd`
cd src
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
#go build
cd ..
[ -d bin ] || mkdir bin
cp -f src/src bin/dsp
[ -d config ] || mkdir config
cp src/config/*.yaml config/
