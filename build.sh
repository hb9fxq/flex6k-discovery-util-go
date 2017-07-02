#!/bin/bash
PATH="/usr/local/go/bin:$PATH"
export GOPATH=$(pwd)

rm -rf bin/*

#pfsense
export GOARCH=amd64
export GOOS=linux
go install github.com/krippendorf/flex6k-discovery-util-go

export GOARCH=amd64
export GOOS=freebsd
go install github.com/krippendorf/flex6k-discovery-util-go

#pfsense
export GOARCH=386
export GOOS=freebsd
go install github.com/krippendorf/flex6k-discovery-util-go

export GOARCH=386
export GOOS=linux
go install github.com/krippendorf/flex6k-discovery-util-go

# raspi
export GOARCH=arm
export GOOS=linux
export GOARM=5

go install github.com/krippendorf/flex6k-discovery-util-go

