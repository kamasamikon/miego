#!/bin/sh -x

CGO_ENABLED=0 go build -ldflags "-w -s" msb.go
upx msb
