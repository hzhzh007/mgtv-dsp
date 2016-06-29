#!/bin/bash

export GOPATH=$GOPATH:`pwd`
cd src
go build
cd ..
[ -d bin ] || mkdir bin
cp -f src/src bin/dsp
[ -d config ] || mkdir config
cp src/config/*.yaml config/
