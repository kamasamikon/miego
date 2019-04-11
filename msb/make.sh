#!/bin/sh

cd main
CGO_ENABLED=0 go build -ldflags "-w -s" main.go
