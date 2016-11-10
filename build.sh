#!/bin/bash

export GOPATH=$(pwd)


export GOARCH=amd64
export GOOS=linux

go install github.com/krippendorf/flex6k-discovery-util-go
mv bin/flex6k-discovery-util-go bin/flex6k-discovery-util-go-linux-amd64

export GOARCH=amd64
export GOOS=freebsd
go install github.com/krippendorf/flex6k-discovery-util-go

export GOARCH=arm
export GOOS=freebsd
export GOARM=7

go install github.com/krippendorf/flex6k-discovery-util-go

