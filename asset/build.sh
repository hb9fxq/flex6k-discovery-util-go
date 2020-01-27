#!/bin/bash
cd ..

# Linux
env GOOS=linux GOARCH=amd64 go build -o ../../../../bin/flex6k-discovery-util-go-build/linux64/flexi
env GOOS=linux GOARCH=386 go build -o ../../../../bin/flex6k-discovery-util-go-build/linux32/flexi

# Raspi
env GOOS=linux GOARCH=arm GOARM=5 go build -o ../../../../bin/flex6k-discovery-util-go-build/raspberryPi/flexi

# Windows
env GOOS=windows GOARCH=amd64 go build -o ../../../../bin/flex6k-discovery-util-go-build/Win64/flexi.exe
env GOOS=windows GOARCH=386 go build -o ../../../../bin/flex6k-discovery-util-go-build/Win32/flexi.exe


# pfsense
env GOOS=freebsd GOARCH=amd64 go build -o ../../../../bin/flex6k-discovery-util-go-build/pfSense64/flexi
env GOOS=freebsd GOARCH=386 go build -o ../../../../bin/flex6k-discovery-util-go-build/pfSense32/flexi
