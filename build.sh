#!/bin/bash
rm -rf build/arm-broadcast build/riscv-broadcast

mkdir "build"
set -e

CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o build/arm-broadcast
CGO_ENABLED=0 GOOS=linux GOARCH=riscv64  go build -o build/riscv-broadcast
