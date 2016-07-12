#!/bin/bash

export GOPATH=$GOPATH:`pwd`
cd src
if [ "$1" == linux ]
then
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
else
    go build
fi
cd ..
[ -d bin ] || mkdir bin
cp -f src/src bin/dsp
[ -d config ] || mkdir config
cp -n  src/config/*.yaml config/
