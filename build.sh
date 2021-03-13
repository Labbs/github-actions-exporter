#!/usr/bin/env bash

set -x

VERSION=`git branch | grep '*' | awk '{print $2}'`
if [[ $VERSION == "master" ]];
then
  VERSION=`git describe --tags --abbrev=0`
fi

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-X 'main.version=$VERSION'"  -o bin/app
