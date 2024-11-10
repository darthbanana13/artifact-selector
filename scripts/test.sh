#!/usr/bin/env bash

if [[ ! -x "$(which go)" ]]; then
    echo "Error! Go executable not found in system path. Please install Go version 1.23 or higher"
    exit 1
fi

go test -v  github.com/darthbanana13/artifact-selector/pkg/log
