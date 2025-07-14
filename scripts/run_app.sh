#!/usr/bin/env bash

if [[ ! -x "$(which go)" ]]; then
    echo "Error! Go executable not found in system path. Please install Go version 1.23 or higher"
    exit 1
fi

# if [[ ! -x "$(which jq)" ]]; then
#   go run cmd/selector/main.go "$@"
#   exit 0
# fi

# go run cmd/selector/main.go "$@" | jq '.'
go run cmd/selector/main.go "$@"
