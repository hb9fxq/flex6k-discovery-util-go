#!/bin/bash
cd ..

# Linux
env GOOS=linux GOARCH=amd64 go build -o ../../../../bin/flex6k-discovery-util-go/linux64/flexi
env GOOS=linux GOARCH=386 go build -o ../../../../bin/flex6k-discovery-util-go/linux32/flexi

# Raspi
env GOOS=linux GOARCH=arm GOARM=5 go build -o ../../../../bin/flex6k-discovery-util-go/raspberryPi/flexi

# Windows
env GOOS=windows GOARCH=amd64 go build -o ../../../../bin/flex6k-discovery-util-go/Win64/flexi.exe
env GOOS=windows GOARCH=386 go build -o ../../../../bin/flex6k-discovery-util-go/Win32/flexi.exe


# pfsense
env GOOS=freebsd GOARCH=amd64 go build -o ../../../../bin/flex6k-discovery-util-go/pfSense64/flexi
env GOOS=freebsd GOARCH=386 go build -o ../../../../bin/flex6k-discovery-util-go/pfSense32/flexi
