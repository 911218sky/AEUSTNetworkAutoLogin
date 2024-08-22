#!/bin/bash

SOURCE_FILE="./src/main.go"
BUILD_FLAGS="-s -w"
UPX_COMPRESSION_FLAGS="-9"

# Linux
# AMD64
GOOS=linux GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o AEUSTNetworkAutoLogin-Linux-x86_64 $SOURCE_FILE && upx $UPX_COMPRESSION_FLAGS AEUSTNetworkAutoLogin-Linux-x86_64
# 386
GOOS=linux GOARCH=386 go build -ldflags "$BUILD_FLAGS" -o AEUSTNetworkAutoLogin-Linux-x86 $SOURCE_FILE && upx $UPX_COMPRESSION_FLAGS AEUSTNetworkAutoLogin-Linux-x86

# macOS
# AMD64
GOOS=darwin GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o AEUSTNetworkAutoLogin-macOS-x86_64 $SOURCE_FILE && upx $UPX_COMPRESSION_FLAGS AEUSTNetworkAutoLogin-macOS-x86_64
# ARM64
GOOS=darwin GOARCH=arm64 go build -ldflags "$BUILD_FLAGS" -o AEUSTNetworkAutoLogin-macOS-arm64 $SOURCE_FILE && upx $UPX_COMPRESSION_FLAGS AEUSTNetworkAutoLogin-macOS-arm64

# Windows
# AMD64
GOOS=windows GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o AEUSTNetworkAutoLogin-windows-x86_64.exe $SOURCE_FILE && upx $UPX_COMPRESSION_FLAGS AEUSTNetworkAutoLogin-windows-x86_64.exe
# 386
GOOS=windows GOARCH=386 go build -ldflags "$BUILD_FLAGS" -o AEUSTNetworkAutoLogin-windows-x86.exe $SOURCE_FILE && upx $UPX_COMPRESSION_FLAGS AEUSTNetworkAutoLogin-windows-x86.exe

echo "All done!"
